package exchange

import "strings"

type Currency struct {
	Symbol string
	Desc   string
}

func (c Currency) String() string {
	return c.Symbol
}

func (c Currency) Eq(c2 Currency) bool {
	return c.Symbol == c2.Symbol
}

// A->B(A兑换为B)
type CurrencyPair struct {
	CurrencyA Currency
	CurrencyB Currency
}

var (
	UNKNOWN = Currency{"UNKNOWN", ""}

	USD  = Currency{"USD", ""}
	USDT = Currency{"USDT", ""}
	PAX  = Currency{"PAX", "https://www.paxos.com"}
	USDC = Currency{"USDC", "https://www.centre.io"}

	BTC  = Currency{"BTC", "https://bitcoin.org"}
	BCH  = Currency{"BCH", ""}
	BSV  = Currency{"BSV", ""}
	LTC  = Currency{"LTC", ""}
	ETH  = Currency{"ETH", ""}
	ETC  = Currency{"ETC", ""}
	EOS  = Currency{"EOS", ""}
	XRP  = Currency{"XRP", ""}
	FT   = Currency{"FT", "FCoin Token"}
	FMEX = Currency{"FMEX", "FCoin Future Token"}

	BTC_USD = CurrencyPair{BTC, USD}
	LTC_USD = CurrencyPair{LTC, USD}
	ETH_USD = CurrencyPair{ETH, USD}
	ETC_USD = CurrencyPair{ETC, USD}
	BCH_USD = CurrencyPair{BCH, USD}
	BSV_USD = CurrencyPair{BSV, USD}
	XRP_USD = CurrencyPair{XRP, USD}
	EOS_USD = CurrencyPair{EOS, USD}

	BTC_USDT  = CurrencyPair{BTC, USDT}
	LTC_USDT  = CurrencyPair{LTC, USDT}
	BCH_USDT  = CurrencyPair{BCH, USDT}
	BSV_USDT  = CurrencyPair{BSV, USDT}
	ETH_USDT  = CurrencyPair{ETH, USDT}
	ETC_USDT  = CurrencyPair{ETC, USDT}
	EOS_USDT  = CurrencyPair{EOS, USDT}
	XRP_USDT  = CurrencyPair{XRP, USDT}
	FT_USDT   = CurrencyPair{FT, USDT}
	FMEX_USDT = CurrencyPair{FMEX, USDT}

	LTC_BTC = CurrencyPair{LTC, BTC}
	ETH_BTC = CurrencyPair{ETH, BTC}
	ETC_BTC = CurrencyPair{ETC, BTC}
	BCH_BTC = CurrencyPair{BCH, BTC}
	BSV_BTC = CurrencyPair{BSV, BTC}
	XRP_BTC = CurrencyPair{XRP, BTC}
	EOS_BTC = CurrencyPair{EOS, BTC}
	FT_BTC  = CurrencyPair{FT, BTC}

	ETC_ETH = CurrencyPair{ETC, ETH}
	LTC_ETH = CurrencyPair{LTC, ETH}
	EOS_ETH = CurrencyPair{EOS, ETH}

	UNKNOWN_PAIR = CurrencyPair{UNKNOWN, UNKNOWN}
)

func (c CurrencyPair) String() string {
	return c.ToSymbol("_")
}

func (c CurrencyPair) Eq(c2 CurrencyPair) bool {
	return c.String() == c2.String()
}

func NewCurrency(symbol, desc string) Currency {
	switch symbol {
	case "usdt", "USDT":
		return USDT
	case "pax", "PAX":
		return PAX
	case "btc", "BTC":
		return BTC
	case "bch", "BCH":
		return BCH
	case "ltc", "LTC":
		return LTC
	case "eos", "EOS":
		return EOS
	case "ft", "FT":
		return FT
	case "fmex", "FMEX":
		return FMEX
	default:
		return Currency{strings.ToUpper(symbol), desc}
	}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, "_")
	if len(currencys) == 2 {
		return CurrencyPair{NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}

func (pair CurrencyPair) AdaptUsdtToUsd() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USDT) {
		CurrencyB = USD
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

func (pair CurrencyPair) AdaptUsdToUsdt() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USD) {
		CurrencyB = USDT
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

//for to symbol lower , Not practical '==' operation method
func (pair CurrencyPair) ToLower() CurrencyPair {
	return CurrencyPair{Currency{strings.ToLower(pair.CurrencyA.Symbol), ""},
		Currency{strings.ToLower(pair.CurrencyB.Symbol), ""}}
}

func (pair CurrencyPair) Reverse() CurrencyPair {
	return CurrencyPair{pair.CurrencyB, pair.CurrencyA}
}
