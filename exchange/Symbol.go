package exchange

import "strings"

type Currency struct {
	Name string
	Desc string
}

func (c Currency) String() string {
	return c.Name
}

func (c Currency) Eq(c2 Currency) bool {
	return c.Name == c2.Name
}

type Symbol struct {
	BaseCurrency  Currency
	QuoteCurrency Currency
}

var (
	UNKNOWN = Currency{"UNKNOWN", ""}

	USDT = Currency{"USDT", ""}
	USD  = Currency{"USD", ""}
	PAX  = Currency{"PAX", "https://www.paxos.com"}

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

	BTC_USD = Symbol{BTC, USD}
	LTC_USD = Symbol{LTC, USD}
	ETH_USD = Symbol{ETH, USD}
	ETC_USD = Symbol{ETC, USD}
	BCH_USD = Symbol{BCH, USD}
	BSV_USD = Symbol{BSV, USD}
	XRP_USD = Symbol{XRP, USD}
	EOS_USD = Symbol{EOS, USD}

	BTC_USDT  = Symbol{BTC, USDT}
	LTC_USDT  = Symbol{LTC, USDT}
	BCH_USDT  = Symbol{BCH, USDT}
	BSV_USDT  = Symbol{BSV, USDT}
	ETH_USDT  = Symbol{ETH, USDT}
	ETC_USDT  = Symbol{ETC, USDT}
	EOS_USDT  = Symbol{EOS, USDT}
	XRP_USDT  = Symbol{XRP, USDT}
	FT_USDT   = Symbol{FT, USDT}
	FMEX_USDT = Symbol{FMEX, USDT}

	LTC_BTC = Symbol{LTC, BTC}
	ETH_BTC = Symbol{ETH, BTC}
	ETC_BTC = Symbol{ETC, BTC}
	BCH_BTC = Symbol{BCH, BTC}
	BSV_BTC = Symbol{BSV, BTC}
	XRP_BTC = Symbol{XRP, BTC}
	EOS_BTC = Symbol{EOS, BTC}
	FT_BTC  = Symbol{FT, BTC}

	ETC_ETH = Symbol{ETC, ETH}
	LTC_ETH = Symbol{LTC, ETH}
	EOS_ETH = Symbol{EOS, ETH}

	UNKNOWN_SYMBOL = Symbol{UNKNOWN, UNKNOWN}
)

func (c Symbol) String() string {
	return c.ToSymbol("/")
}

func (c Symbol) Eq(c2 Symbol) bool {
	return c.String() == c2.String()
}

func NewCurrency(name, desc string) Currency {
	switch name {
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
		return Currency{strings.ToUpper(name), desc}
	}
}

func NewSymbol(baseCurrency Currency, quoteCurrency Currency) Symbol {
	return Symbol{baseCurrency, quoteCurrency}
}

func NewSymbol2(symbolName string) Symbol {
	currencys := strings.Split(symbolName, "_")
	if len(currencys) == 2 {
		return Symbol{
			NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], ""),
		}
	}

	return UNKNOWN_SYMBOL
}

func (symbol Symbol) ToSymbol(joinChar string) string {
	return strings.Join([]string{symbol.BaseCurrency.Name, symbol.QuoteCurrency.Name}, joinChar)
}

func (symbol Symbol) ToSymbol2(joinChar string) string {
	return strings.Join([]string{symbol.QuoteCurrency.Name, symbol.BaseCurrency.Name}, joinChar)
}

func (symbol Symbol) AdaptUsdtToUsd() Symbol {
	quoteCurrency := symbol.QuoteCurrency
	if symbol.QuoteCurrency.Eq(USDT) {
		quoteCurrency = USD
	}
	return Symbol{symbol.BaseCurrency, quoteCurrency}
}

func (symbol Symbol) AdaptUsdToUsdt() Symbol {
	quoteCurrency := symbol.QuoteCurrency
	if symbol.QuoteCurrency.Eq(USD) {
		quoteCurrency = USDT
	}
	return Symbol{symbol.BaseCurrency, quoteCurrency}
}

//for to symbol lower, Not practical '==' operation method
func (symbol Symbol) ToLower() Symbol {
	return Symbol{
		Currency{strings.ToLower(symbol.BaseCurrency.Name), ""},
		Currency{strings.ToLower(symbol.QuoteCurrency.Name), ""},
	}
}

func (symbol Symbol) Reverse() Symbol {
	return Symbol{symbol.QuoteCurrency, symbol.BaseCurrency}
}
