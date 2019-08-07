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
	PAX  = Currency{"PAX", ""}
	TUSD = Currency{"TUSD", ""}
	BTC  = Currency{"BTC", ""}
	BCH  = Currency{"BCH", ""}
	BSV  = Currency{"BSV", ""}
	ETH  = Currency{"ETH", ""}
	ETC  = Currency{"ETC", ""}
	LTC  = Currency{"LTC", ""}
	EOS  = Currency{"EOS", ""}
	XRP  = Currency{"XRP", ""}
	FT   = Currency{"FT", "FCoin Token"}
	FMEX = Currency{"FMEX", "FCoin Contract Token"}

	PAX_USDT  = Symbol{PAX, USDT}
	TUSD_USDT = Symbol{TUSD, USDT}
	BTC_USDT  = Symbol{BTC, USDT}
	BCH_USDT  = Symbol{BCH, USDT}
	BSV_USDT  = Symbol{BSV, USDT}
	ETH_USDT  = Symbol{ETH, USDT}
	ETC_USDT  = Symbol{ETC, USDT}
	LTC_USDT  = Symbol{LTC, USDT}
	EOS_USDT  = Symbol{EOS, USDT}
	XRP_USDT  = Symbol{XRP, USDT}
	FT_USDT   = Symbol{FT, USDT}
	FMEX_USDT = Symbol{FMEX, USDT}

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
	case "tusd", "TUSD":
		return TUSD
	case "btc", "BTC":
		return BTC
	case "bch", "BCH":
		return BCH
	case "bsv", "BSV":
		return BSV
	case "eth", "ETH":
		return ETH
	case "etc", "ETC":
		return ETC
	case "ltc", "LTC":
		return LTC
	case "eos", "EOS":
		return EOS
	case "xrp", "XRP":
		return XRP
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
