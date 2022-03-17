package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func NewParserService() *ParserService {
	return &ParserService{}
}

func (p *ParserService) ParseCSV() ([][]string, error) {
	file, err := os.Open(INPUT_PATH)
	if err != nil {
		fmt.Println("file open error")
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	orderBookList := p.parseOrderBookList(scanner)

	return orderBookList, nil
}

func (p *ParserService) parseOrderBookList(scanner *bufio.Scanner) [][]string {
	output := [][]string{}
	currentBookInput := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		} else if line[0:1] == FLUSH_ORDER_BOOK {
			currentBookInput = append(currentBookInput, line)
			output = append(output, currentBookInput)
			currentBookInput = []string{}
		} else if line[0:1] == CANCEL_ORDER || line[0:1] == NEW_ORDER {
			currentBookInput = append(currentBookInput, line)
		}
	}
	return output
}

// TransformOrderBookListData: converts raw string data to structured format
func (p *ParserService) TransformOrderBookListData(orderBookList [][]string) ([][]Order, error) {
	var (
		orderBookInputs [][]Order
		orderBook       []Order
		order           Order
	)

	for _, orderBookInput := range orderBookList {
		for _, orderline := range orderBookInput {
			orderSplit := strings.Split(orderline, ",")
			command := strings.TrimSpace(orderSplit[0])
			if command == FLUSH_ORDER_BOOK {
				if len(orderSplit) != 1 {
					return nil, errors.New("FLUSH_ORDER_BOOK invalid input in TransformOrderBookListData()")
				}
				order = Order{
					Command: FLUSH_ORDER_BOOK,
				}
			} else if command == NEW_ORDER {
				if len(orderSplit) != 7 {
					return nil, errors.New("NEW_ORDER invalid input in TransformOrderBookListData")
				}
				userID, err := strconv.Atoi(strings.TrimSpace(orderSplit[1]))
				if err != nil {
					err = errors.Wrap(err, "error converting UserID in NEW_ORDER for TransformOrderBookListData()")
					return nil, err
				}

				price, err := strconv.Atoi(strings.TrimSpace(orderSplit[3]))
				if err != nil {
					err = errors.Wrap(err, "error converting price in NEW_ORDER for TransformOrderBookListData()")
					return nil, err
				}

				quantity, err := strconv.Atoi(strings.TrimSpace(orderSplit[4]))
				if err != nil {
					err = errors.Wrap(err, "error converting quantity in NEW_ORDER for TransformOrderBookListData()")
					return nil, err
				}
				userOrderID, err := strconv.Atoi(strings.TrimSpace(orderSplit[6]))
				if err != nil {
					err = errors.Wrap(err, "error converting userOrderID in NEW_ORDER for TransformOrderBookListData()")
					return nil, err
				}

				order = Order{
					Command:     NEW_ORDER,
					UserID:      userID,
					Symbol:      strings.TrimSpace(orderSplit[2]),
					Price:       price,
					Quantity:    quantity,
					Side:        strings.TrimSpace(orderSplit[5]),
					UserOrderID: userOrderID,
				}
			} else if command == CANCEL_ORDER {
				if len(orderSplit) != 3 {
					return nil, errors.New("CANCEL_ORDER invalid input in TransformOrderBookListData()")
				}
				userID, err := strconv.Atoi(strings.TrimSpace(orderSplit[1]))
				if err != nil {
					err = errors.Wrap(err, "error converting userID in CANCEL_ORDER for TransformOrderBookListData()")
					return nil, err
				}

				userOrderID, err := strconv.Atoi(strings.TrimSpace(orderSplit[2]))
				if err != nil {
					err = errors.Wrap(err, "error converting userOrderID in CANCEL_ORDER for TransformOrderBookListData()")
					return nil, err
				}

				order = Order{
					Command:     CANCEL_ORDER,
					UserID:      userID,
					UserOrderID: userOrderID,
				}
			}
			orderBook = append(orderBook, order)
		}
		orderBookInputs = append(orderBookInputs, orderBook)
		orderBook = []Order{}
	}
	return orderBookInputs, nil
}
