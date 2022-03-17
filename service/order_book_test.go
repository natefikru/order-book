package service

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func compareErrors(expected error, returned error) bool {
	if expected == nil {
		return returned == nil
	}
	return reflect.TypeOf(expected) == reflect.TypeOf(returned) && expected.Error() == returned.Error()
}

func TestNewOrder(t *testing.T) {
	testService := NewOrderBookService()
	tests := map[string]struct {
		order   *Order
		topBook TopBook
		output  string
		err     error
	}{
		"Valid Buy Order": {
			order: &Order{
				UserID:      1,
				UserOrderID: 1,
				Command:     NEW_ORDER,
				Symbol:      "IBM",
				Price:       10,
				Quantity:    100,
				Side:        BUY,
			},
			topBook: TopBook{
				Price: 20,
			},
			err:    nil,
			output: "A, 1, 1",
		},
		"Valid Sell Order": {
			order: &Order{
				UserID:      1,
				UserOrderID: 2,
				Command:     NEW_ORDER,
				Symbol:      "IBM",
				Price:       10,
				Quantity:    100,
				Side:        SELL,
			},
			topBook: TopBook{
				Price: 5,
			},
			err:    nil,
			output: "A, 1, 2",
		},
		"Invalid Order Price": {
			order: &Order{
				UserID:      1,
				UserOrderID: 3,
				Command:     NEW_ORDER,
				Symbol:      "IBM",
				Price:       0,
				Quantity:    100,
				Side:        SELL,
			},
			err:    errors.New("Error creating new order for order 3: invalid price 0"),
			output: "",
		},
		"Rejected Buy Order": {
			order: &Order{
				UserID:      2,
				UserOrderID: 4,
				Command:     NEW_ORDER,
				Symbol:      "IBM",
				Price:       10,
				Quantity:    100,
				Side:        BUY,
			},
			topBook: TopBook{
				Price: 5,
			},
			err:    nil,
			output: "R, 2, 4",
		},
		"Rejected Sell Order": {
			order: &Order{
				UserID:      3,
				UserOrderID: 5,
				Command:     NEW_ORDER,
				Symbol:      "IBM",
				Price:       10,
				Quantity:    100,
				Side:        SELL,
			},
			topBook: TopBook{
				Price: 20,
			},
			err:    nil,
			output: "R, 3, 5",
		},
	}

	for name, test := range tests {
		testService.OrderBook = NewOrderBook()

		if test.order.Side == BUY {
			testService.OrderBook.TopBookAsk = test.topBook
		} else if test.order.Side == SELL {
			testService.OrderBook.TopBookBid = test.topBook
		}
		output, err := testService.newOrder(test.order)
		if !compareErrors(test.err, err) {
			t.Errorf("Expected error %s, received %s for test %s", test.err, err, name)
		}
		if output != test.output {
			t.Errorf("Expected output %s, received %s for test %s", test.output, output, name)
		}
	}
}

func TestCancelOrder(t *testing.T) {
	testService := NewOrderBookService()
	tests := map[string]struct {
		order          *Order
		existingOrders []Order
		finalAskLength int
		finalBidLength int
		output         string
		err            error
	}{
		"Valid cancel order for bid": {
			order: &Order{
				Side:        BUY,
				UserID:      1,
				UserOrderID: 1,
			},
			existingOrders: []Order{
				{
					Side:        BUY,
					UserID:      1,
					UserOrderID: 1,
				},
			},
			output: "A, 1, 1",
			err:    nil,
		},
		"Valid cancel order for bid 2": {
			order: &Order{
				Side:        BUY,
				UserID:      13,
				UserOrderID: 10,
			},
			existingOrders: []Order{
				{
					Side:        BUY,
					UserID:      13,
					UserOrderID: 10,
				},
				{
					Side:        BUY,
					UserID:      13,
					UserOrderID: 12,
				},
			},
			finalBidLength: 1,
			output:         "A, 13, 10",
			err:            nil,
		},
		"Valid cancel order for ask": {
			order: &Order{
				Side:        SELL,
				UserID:      2,
				UserOrderID: 2,
			},
			existingOrders: []Order{
				{
					Side:        SELL,
					UserID:      2,
					UserOrderID: 2,
				},
			},
			output: "A, 2, 2",
			err:    nil,
		},
		"Valid cancel order for ask 2": {
			order: &Order{
				Side:        SELL,
				UserID:      15,
				UserOrderID: 15,
			},
			existingOrders: []Order{
				{
					Side:        SELL,
					UserID:      15,
					UserOrderID: 15,
				},
				{
					Side:        SELL,
					UserID:      16,
					UserOrderID: 16,
				},
			},
			finalAskLength: 1,
			output:         "A, 15, 15",
			err:            nil,
		},
		"Invalid order side": {
			order: &Order{
				Side:        "K",
				UserID:      2,
				UserOrderID: 2,
			},
			existingOrders: []Order{
				{
					Side:        BUY,
					UserID:      14,
					UserOrderID: 14,
				},
			},
			finalBidLength: 1,
			output:         "",
			err:            nil,
		},
	}

	for name, test := range tests {
		testService.OrderBook = NewOrderBook()

		for _, existingOrder := range test.existingOrders {
			if existingOrder.Side == BUY {
				testService.OrderBook.Bids = append(testService.OrderBook.Bids, existingOrder)
			} else if existingOrder.Side == SELL {
				testService.OrderBook.Asks = append(testService.OrderBook.Asks, existingOrder)
			}
		}
		output, err := testService.cancelOrder(test.order)
		if !compareErrors(test.err, err) {
			t.Errorf("Expected error %s, received %s for test %s", test.err, err, name)
		}
		if output != test.output {
			t.Errorf("Expected output %s, received %s for test %s", test.output, output, name)
		}
		if test.finalAskLength != len(testService.OrderBook.Asks) {
			t.Errorf("Expected ask length %v, received %v for test %s", test.finalAskLength, len(testService.OrderBook.Asks), name)
		}
		if test.finalBidLength != len(testService.OrderBook.Bids) {
			t.Errorf("Expected bid length %v, received %v for test %s", test.finalBidLength, len(testService.OrderBook.Bids), name)
		}
	}
}
