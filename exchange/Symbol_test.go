package exchange

import (
	"testing"
)

func TestCurrency2_String(t *testing.T) {
	btc := NewCurrency("btc", "bitcoin")
	ltc := NewCurrency("ltc", "litecoin")
	btc2 := Currency{"BTC", "bitcoin.org"}

	t.Log(btc == BTC)
	t.Log(ltc.Desc, btc.Desc)
	t.Log(btc == btc2)
}

func TestSymbol2_String(t *testing.T) {
	btc_usdt := NewSymbol(NewCurrency("btc", ""), NewCurrency("usdt", ""))

	t.Log(btc_usdt.String() == "BTC_USDT")
	t.Log(btc_usdt.ToLower().ToSymbol("") == "btcusdt")
	t.Log(btc_usdt.ToLower().String() == "btc_usdt")
	t.Log(btc_usdt.Reverse().String() == "USDT_BTC")
	t.Log(btc_usdt.Eq(BTC_USDT))
}
