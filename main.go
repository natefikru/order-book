package main

import (
	"fmt"
	"order_book_exercise/service"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("Order Book Excercise Started")
	// Parse raw data
	parserService := service.NewParserService()
	orderBookListData, err := parserService.ParseCSV()
	if err != nil {
		err = errors.Wrap(err, "error parsing CSV in main function")
		fmt.Println(err.Error())
		panic(err)
	}
	// Transform raw data to supported struct
	orderBooks, err := parserService.TransformOrderBookListData(orderBookListData)
	if err != nil {
		err = errors.Wrap(err, "error transforming raw string data in main function")
		fmt.Println(err.Error())
		panic(err)
	}

	orderbookService := service.NewOrderBookService()
	// iterate through each constructed orderbook order, process and log the results
	// TODO: Execute orderbook processing via threaded go routines and save the data to external database.
	for i, orderBook := range orderBooks {

		// Create a fresh orderbook and attach to the service
		orderbookService.OrderBook = service.NewOrderBook()

		fmt.Printf("Processing Order book %v\n", i+1)
		err = orderbookService.ProcessOrderBook(orderBook)
		if err != nil {
			err = errors.Wrapf(err, "error Processing order book %v in main function", i+1)
			fmt.Println(err.Error())
			panic(err)
		}
		fmt.Println()
	}
}
