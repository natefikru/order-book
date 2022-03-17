# Order Book Programming Exercise
#### Nathnael Fikru 


## How to Run
The project is written in Go, you will need to install this project into your local computer's `$GOPATH/src` directory.Navigate to this directory via the command line. Once in the directory, run the `go get` command to install the neccesary dependencies.

Before you run there are a few things to know about the projects configuration. If you open up the `/service/service.go` file, you will see a list of consts. 

`INPUT_PATH` will determine which input file to read from
`IS_TRADING_ENABLED` will determine whether or not we want to turn on the optional trading mode. This value is set to `false` by default.

To run the program, navigate to the root folder and run the command
```
go run main.go
```
## How to Run Tests
All of the tests for this project have been written within the `/service` directory. To run the tests, navigate to that folder and run 
```
go test
```
This will execute the test suite.

## Overview
Overall this challenge was a fun one to tackle. there are definitely things that I wish I was able to implement in code but given the 24 hour time constraint, certain things were not possible.

### Project structure
For the function section of the project, I tried breaking it down into 2 main pieces. The `ParserService` that would read the input of the text file specified, The `OrderBookService` that would do most of the heavy lifting for interpereting the data. The top level `main` function invokes both of these services.

#### ParserService
The parser is a straightforward service that essentially ignores all lines that don't begin with the reserved commands that we are looking for within the OrderBookService. Once a valid line is found, we parse that line based on its command type.

once all of the data has been parsed, we then transform the raw text data to structured orderBook objects for our orderbook service to handle

#### OrderBookService
At the top level, the OrderBookServices holds the `IsTradingEnabled` configuration as well as a fresh OrderBook per scenario. 
The OrderBook struct has the following Attributes associated with it.
```
type OrderBook struct {
	Bids []Order // slice of bid orders in order from highest to lowest
	Asks []Order // slice of ask orders in order from lowest to highest

	TopBookBid TopBook // location of topBook bid  data
	TopBookAsk TopBook // location of topBook bid  data

	OrderDict map[int]string // hashmap that stores order side information for easy lookup
}
```

As we process each order that comes along, this data in the handler gets updated as it goes. The `TopBook` data needs to be reassessed after every single new order and cancel order command to determine if anything has changed.

#### Tests
The tests have been created using the go test framework. I chose to go with table formatted tests to reduce the amount of redundant code within the testing file. There are tests for `newOrder` and `cancelOrder`. 

### Things to do if I had more time
Theres a number of things that i would do to make this service more robust.
- First of all, to simplify the application installationg and run process, I would dockerize the project and print the output out to a file instead of the standard output. 
- I would write more tests that encompass each of the core commands including executeTrade() and handleTopOfBook(). In additon I would create more scenarios. It would also be nice if I was able to write a test that took in both the input.csv file and output.csv file and cross checked them against eachother. There should also be tests on the parser as well.
- Instead of each core command passing back strings back to the ProcessOrderBook function to be printed, I would create a new data structure specifically for the data that is being outputted. This would make it easier to handle the output data and give me more flexibility with the output format as well as flexibility in testing.
- I created a simple executeTrade() function that essentially just looks at the topBook order but, if I had more time I would build out the functionality to be able to split a trade over multiple orders within the queue and break down over multiple quantities.
- I chose to run the orderbookservice synchronously choosing to avoid the possiblity of different orderBook scenario standard outputs overlapping eachother in the wrong order. There are definitely many reasons to utilize goroutines here to process many orderbooks at the same time. I didn't see a solid use case here though since the data is staying static and the load was very low.

### Time/Space Complexity
**newOrder = Time: O(n) Space: O(n)** - must loop through bids or asks to find insertionIndex, creates a new OrderList to replace the old one, 
**cancelOrder = Time: O(n) Space: O(n)** - must loop through bids or asks to find deletionIndex, creates a new OrderList to replace the old one, 
**execute Trade Time:O(1) Space O(n)**  = Only looks at top values, and  creates a new OrderList to replace the old one.  
**evaluateBook = Time: O(n) Space: O(1)** - must loop through bids or asks to evaulate the topBook Value. Only Creates 1 new topBook Value