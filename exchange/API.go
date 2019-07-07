package exchange

// api interface

type API interface {
	GetExchangeName() string

	GetAccount() (*Account, error)
	GetSubAccount(currency Currency) (*SubAccount, error)

	LimitBuy(amount, price string, currency Symbol) (*Order, error)
	LimitSell(amount, price string, currency Symbol) (*Order, error)
	MarketBuy(amount, price string, currency Symbol) (*Order, error)
	MarketSell(amount, price string, currency Symbol) (*Order, error)

	CancelOrder(orderId string, currency Symbol) (bool, error)
	GetOrder(orderId string, currency Symbol) (*NewOrder, error)
	GetActiveOrders(currency Symbol) ([]NewOrder, error)
	GetOrderHistorys(currency Symbol, currentPage, pageSize int) ([]NewOrder, error)

	GetTicker(currency Symbol) (*Ticker, error)
	GetDepth(size int, currency Symbol) (*Depth, error)
	GetKlineRecords(currency Symbol, period, size, since int) ([]Kline, error)

	GetTrades(symbol Symbol, since int64) ([]Trade, error)
}
