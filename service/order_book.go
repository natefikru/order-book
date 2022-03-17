package service

import (
	"fmt"

	"github.com/pkg/errors"
)

/////////////////////////
///   INITIALIZERS   ////
/////////////////////////

func NewOrderBookService() *OrderBookService {
	return &OrderBookService{
		IsTradingEnabled: IS_TRADING_ENABLED,
	}
}

// NewOrderBook: Initializes Order Book within OrderBookService
func NewOrderBook() OrderBook {
	return OrderBook{
		TopBookBid: TopBook{
			UserID: 0,
			Price:  MIN_INT,
		},
		TopBookAsk: TopBook{
			UserID: 0,
			Price:  MAX_INT,
		},
		OrderDict: make(map[int]string),
	}
}

/////////////////////////
///       MAIN       ////
/////////////////////////

// ProcessOrderBook: Main function processes order book limit bids/asks by price and time
func (o *OrderBookService) ProcessOrderBook(orderBook []Order) error {
	var output string
	var err error
	for _, order := range orderBook {
		switch order.Command {
		case NEW_ORDER:
			output, err = o.newOrder(&order)
			if err != nil {
				return errors.Wrapf(err, "error creating new order in ProcessOrderBook for order: %v", order.UserOrderID)
			}
			if output != "" {
				fmt.Println(output)
			}

			// Based on configuration.
			if o.IsTradingEnabled {
				output, err = o.executeTrade()
				if err != nil {
					return errors.Wrapf(err, "error attempting to execute a trade in ProcessOrderBook for order: %v", order.UserOrderID)
				}
				if output != "" {
					fmt.Println(output)
				}
			}

			output, err = o.handleTopOfBook()
			if err != nil {
				return errors.Wrapf(err, "error handling top of book for new order in ProcessOrderBook for order: %v", order.UserOrderID)
			}
			if output != "" {
				fmt.Println(output)
			}
		case CANCEL_ORDER:
			// Instead of searching through both asks and bids list to find which order to cancel, we keep track of
			// what side an order is on via an in memory hashmap o.OrderBook.OrderDict
			order.Side = o.OrderBook.OrderDict[order.UserOrderID]
			output, err := o.cancelOrder(&order)
			if err != nil {
				return errors.Wrapf(err, "error cancelling order in ProcessOrderBook for order: %v", order.UserOrderID)
			}
			if output != "" {
				fmt.Println(output)
			}
			output, err = o.handleTopOfBook()
			if err != nil {
				return errors.Wrapf(err, "error handling top of book for cancel order in ProcessOrderBook for order: %v", order.UserOrderID)
			}
			if output != "" {
				fmt.Println(output)
			}
		case FLUSH_ORDER_BOOK:
			o.flushBook()
		}
	}
	return nil
}

/////////////////////////
///   CORE COMMANDS  ////
/////////////////////////

// if TRADING_IS_ENABLED == True, the program executes cross book orders as trades
// TODO: Support trading multiple orders instead of just lowestAsk and highest Bid.
func (o *OrderBookService) executeTrade() (string, error) {
	var output string
	if len(o.OrderBook.Asks) != 0 && len(o.OrderBook.Bids) != 0 {
		lowestAsk := o.OrderBook.Asks[0]
		highestBid := o.OrderBook.Bids[0]
		if highestBid.Price >= lowestAsk.Price && highestBid.Quantity == lowestAsk.Quantity {
			orderList, err := remove(o.OrderBook.Asks, 0)
			if err != nil {
				errors.Wrap(err, "error removing ask in trade call for newOrder()")
			}
			o.OrderBook.Asks = orderList

			orderList, err = remove(o.OrderBook.Bids, 0)
			if err != nil {
				errors.Wrap(err, "error removing bid in trade call for newOrder()")
			}
			o.OrderBook.Bids = orderList

			output = fmt.Sprintf("T, %v, %v, %v %v, %v, %v", highestBid.UserID, highestBid.UserOrderID, lowestAsk.UserID, lowestAsk.UserOrderID, lowestAsk.Price, lowestAsk.Quantity)
			return output, nil
		}
	}
	return "", nil
}

// newOrder: function that creates a brand new order within the order book
// also evaluates orders that attempt to cross the book
func (o *OrderBookService) newOrder(order *Order) (string, error) {
	var output string
	if order.Price == 0 {
		return "", errors.Errorf("Error creating new order for order %v: invalid price %v", order.UserOrderID, order.Price)
	}
	if order.Side == BUY {
		if !o.IsTradingEnabled && order.Price >= o.OrderBook.TopBookAsk.Price {
			output = fmt.Sprintf("R, %v, %v", order.UserID, order.UserOrderID)
			return output, nil
		}
		insertionIndex := len(o.OrderBook.Bids)
		for i, bidOrder := range o.OrderBook.Bids {
			if order.Price > bidOrder.Price {
				insertionIndex = i
				break
			}
		}
		bids, err := insertOrder(o.OrderBook.Bids, insertionIndex, *order)
		if err != nil {
			return "", errors.Wrap(err, "error inserting order to bids in newOrder()")
		}
		o.OrderBook.Bids = bids
		o.OrderBook.OrderDict[order.UserOrderID] = BUY
		output = fmt.Sprintf("A, %v, %v", order.UserID, order.UserOrderID)
		return output, nil
	} else if order.Side == SELL {
		if !o.IsTradingEnabled && order.Price <= o.OrderBook.TopBookBid.Price {
			output = fmt.Sprintf("R, %v, %v", order.UserID, order.UserOrderID)
			return output, nil
		}
		insertionIndex := len(o.OrderBook.Asks)
		for i, askOrder := range o.OrderBook.Asks {
			if order.Price < askOrder.Price {
				insertionIndex = i
				break
			}
		}
		asks, err := insertOrder(o.OrderBook.Asks, insertionIndex, *order)
		if err != nil {
			return "", errors.Wrap(err, "error inserting order to asks in newOrder()")
		}
		o.OrderBook.Asks = asks
		o.OrderBook.OrderDict[order.UserOrderID] = SELL
		output = fmt.Sprintf("A, %v, %v", order.UserID, order.UserOrderID)
		return output, nil
	}
	return "", nil
}

// cancelOrder: cancels orders within the orderbook by ID
func (o *OrderBookService) cancelOrder(order *Order) (string, error) {
	var output string
	if order.Side == BUY {
		bids := o.OrderBook.Bids
		for i := range bids {
			if bids[i].UserOrderID == order.UserOrderID {
				orderList, err := remove(o.OrderBook.Bids, i)
				if err != nil {
					return "", errors.Wrap(err, "error removing bid in cancelOrder()")
				}
				o.OrderBook.Bids = orderList
				output = fmt.Sprintf("A, %v, %v", order.UserID, order.UserOrderID)
				return output, nil
			}
		}
	} else if order.Side == SELL {
		asks := o.OrderBook.Asks
		for i := range asks {
			if asks[i].UserOrderID == order.UserOrderID {
				orderList, err := remove(o.OrderBook.Asks, i)
				if err != nil {
					return "", errors.Wrap(err, "error removing ask in cancelOrder()")
				}
				o.OrderBook.Asks = orderList
				output = fmt.Sprintf("A, %v, %v", order.UserID, order.UserOrderID)
				return output, nil
			}
		}
	}
	return "", nil
}

// flushBook: clears the orderbook
func (o *OrderBookService) flushBook() {
	o.OrderBook = OrderBook{}
}

/////////////////////////
///   TOP OF BOOK   ////
////////////////////////

// handleTopOfBook: Determines if we need to handle the top of book for asks or bids
func (o *OrderBookService) handleTopOfBook() (string, error) {
	output, err := o.evaluateBook(BUY, o.OrderBook.TopBookBid, o.OrderBook.Bids)
	if err != nil {
		return "", errors.Wrap(err, "error for bid order in assessTopOfBook()")
	}
	if output != "" {
		return output, nil
	}
	output, err = o.evaluateBook(SELL, o.OrderBook.TopBookAsk, o.OrderBook.Asks)
	if err != nil {
		return "", errors.Wrap(err, "error for ask order in assessTopOfBook()")
	}
	if output != "" {
		return output, nil
	}
	return "", nil
}

// evaluateBook: groups top orders that share UserID's and Price, thenwe add the quantities together and store to handler
func (o *OrderBookService) evaluateBook(side string, currentTopOfBook TopBook, orders []Order) (string, error) {
	var output string
	if side != BUY && side != SELL {
		return "", errors.New("invalid side input in evaluatebook")
	}
	var newTopOfBook TopBook
	if side == BUY {
		newTopOfBook = TopBook{
			UserID:   0,
			Price:    MIN_INT,
			Quantity: 0,
		}
	} else if side == SELL {
		newTopOfBook = TopBook{
			UserID:   0,
			Price:    MAX_INT,
			Quantity: 0,
		}
	}

	if len(orders) != 0 {
		newTopOfBook.UserID = orders[0].UserID
		newTopOfBook.Price = orders[0].Price
		newTopOfBook.Quantity = orders[0].Quantity

		for i := 1; i < len(orders); i++ {
			if orders[i].UserID == newTopOfBook.UserID && orders[i].Price == newTopOfBook.Price {
				newTopOfBook.Quantity += orders[i].Quantity
			}
		}
	}
	if newTopOfBook.UserID != currentTopOfBook.UserID ||
		newTopOfBook.Quantity != currentTopOfBook.Quantity ||
		newTopOfBook.Price != currentTopOfBook.Price {

		if side == BUY {
			o.OrderBook.TopBookBid = newTopOfBook
			if newTopOfBook.Price == MIN_INT {
				output = "B, B, -, -"
			} else {
				output = fmt.Sprintf("B, B, %v, %v", newTopOfBook.Price, newTopOfBook.Quantity)
			}
		} else if side == SELL {
			o.OrderBook.TopBookAsk = newTopOfBook
			if newTopOfBook.Price == MAX_INT {
				output = "B, S, -, -"
			} else {
				output = fmt.Sprintf("B, S, %v, %v", newTopOfBook.Price, newTopOfBook.Quantity)
			}
		}
	}
	return output, nil
}

/////////////////////////
///    HELPERS      ////
////////////////////////
func insertOrder(orderList []Order, index int, newOrder Order) ([]Order, error) {
	if index < 0 || index > len(orderList) {
		return nil, errors.New("index out of bounds in insertOrder()")
	}
	if len(orderList) == index {
		return append(orderList, newOrder), nil
	}
	orderList = append(orderList[:index+1], orderList[index:]...)
	orderList[index] = newOrder
	return orderList, nil
}

func remove(orderList []Order, index int) ([]Order, error) {
	if index < 0 || index >= len(orderList) {
		return nil, errors.New("index out of bounds in remove()")
	}
	orderList = append(orderList[:index], orderList[index+1:]...)
	return orderList, nil
}
