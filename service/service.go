package service

const (
	// CONFIGURATION
	INPUT_PATH         = "input_file.csv"
	IS_TRADING_ENABLED = false

	// RESERVED COMMANDS AND SIDE SIGNIFIERS
	NEW_ORDER        = "N"
	CANCEL_ORDER     = "C"
	FLUSH_ORDER_BOOK = "F"
	BUY              = "B"
	SELL             = "S"

	UINT_SIZE = 32 << (^uint(0) >> 32 & 1)
	MAX_INT   = 1<<(UINT_SIZE-1) - 1
	MIN_INT   = -MAX_INT - 1
)

type OrderBook struct {
	Bids []Order
	Asks []Order

	TopBookBid TopBook
	TopBookAsk TopBook

	OrderDict map[int]string
}

type TopBook struct {
	UserID   int
	Price    int
	Quantity int
}

type Order struct {
	UserID      int
	UserOrderID int
	Command     string
	Symbol      string
	Price       int
	Quantity    int
	Side        string
}

type OrderBookService struct {
	IsTradingEnabled bool
	OrderBook        OrderBook
}

type ParserService struct {
}
