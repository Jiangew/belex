package exchange

type API interface {
	GetExchangeName() string

	GetAccount() (*Account, error)
	GetSubAccount(currency Currency) (*SubAccount, error)

	LimitBuy(amount, price string, symbol Symbol) (*NewOrder, error)
	LimitSell(amount, price string, symbol Symbol) (*NewOrder, error)

	MarketBuy(amount, price string, symbol Symbol) (*NewOrder, error)
	MarketSell(amount, price string, symbol Symbol) (*NewOrder, error)

	GetActiveOrders(symbol Symbol) ([]NewOrder, error)
	CancelOrder(orderId string, symbol Symbol) (bool, error)

	GetOrder(orderId string, symbol Symbol) (*NewOrder, error)
	GetOrderHistorys(symbol Symbol, currentPage, pageSize int) ([]NewOrder, error)

	GetTicker(symbol Symbol) (*Ticker, error)
	GetDepth(size int, symbol Symbol) (*Depth, error)

	GetKlines(symbol Symbol) ([]Kline, error)
	IsOrderable(symbol Symbol) (bool, error)

	GetTrades(symbol Symbol, since int64) ([]Trade, error)
}
