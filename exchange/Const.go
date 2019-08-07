package exchange

import "fmt"

type TradeSide int

const (
	BUY = 1 + iota
	SELL
)

func (ts TradeSide) String() string {
	switch ts {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

type OrderType int

const (
	LIMIT = 1 + iota
	MARKET
	FOK
	IOC
	POST_ONLY
)

func (ot OrderType) String() string {
	if ot > 0 && int(ot) <= len(orderTypes) {
		return orderTypes[ot-1]
	}
	return fmt.Sprintf("UNKNOWN_ORDER_TYPE(%d)", ot)
}

var orderTypes = [...]string{"LIMIT", "MARKET", "FOK", "IOC", "POST_ONLY"}

type OrderState int

const (
	SUBMITTED = iota
	PARTIAL_FILLED
	PARTIAL_CANCELED
	FILLED
	CANCELED
	PENDING_CANCEL
)

func (state OrderState) String() string {
	return orderStates[state]
}

var orderStates = [...]string{"SUBMITTED", "PARTIAL_FILLED", "PARTIAL_CANCELED", "FILLED", "CANCELED", "PENDING_CANCEL"}

const (
	KLINE_PERIOD_1MIN = 1 + iota
	KLINE_PERIOD_3MIN
	KLINE_PERIOD_5MIN
	KLINE_PERIOD_15MIN
	KLINE_PERIOD_30MIN
	KLINE_PERIOD_60MIN
	KLINE_PERIOD_1H
	KLINE_PERIOD_2H
	KLINE_PERIOD_4H
	KLINE_PERIOD_6H
	KLINE_PERIOD_8H
	KLINE_PERIOD_12H
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_3DAY
	KLINE_PERIOD_1WEEK
	KLINE_PERIOD_1MONTH
	KLINE_PERIOD_1YEAR
)

const (
	FCOIN     = "fcoin.com"
	FCOIN_PRO = "fcoin.pro"
)
