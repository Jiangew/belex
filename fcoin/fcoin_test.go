package fcoin

import (
	"github.com/jiangew/belex/exchange"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var ft = NewFCoin(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1086")
			return nil, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
}, "1412ac27e3f741c796f7c4600069d9f1", "4843754749be46919d986142917f06d7")

func TestFCoin_GetTicker(t *testing.T) {
	t.Log(ft.GetTicker(exchange.NewSymbol2("BTC_USDT")))
}

func TestFCoin_GetDepth(t *testing.T) {
	dep, _ := ft.GetDepth(1, exchange.BTC_USDT)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}

func TestFCoin_GetAccount(t *testing.T) {
	acc, _ := ft.GetAccount()
	t.Log(acc)
}

func TestFCoin_LimitBuy(t *testing.T) {
	t.Log(ft.LimitBuy("1", "0.16", exchange.FMEX_USDT))
}

func TestFCoin_LimitSell(t *testing.T) {
	t.Log(ft.LimitSell("1", "0.23", exchange.FMEX_USDT))
}

func TestFCoin_CancelOrder(t *testing.T) {
	t.Log(ft.CancelOrder("4-AFlOjmvn52YWAgb-LZJ5jaqysio1Icl67nDtt08pYQGpkVbVQT5gYQakM1deT9G4yxp70zjD_IDdU5d6DwgA==", exchange.FMEX_USDT))
}

func TestFCoin_GetUnfinishOrders(t *testing.T) {
	t.Log(ft.GetActiveOrders(exchange.FMEX_USDT))
}

func TestFCoin_GetOrderHistorys(t *testing.T) {
	t.Log(ft.GetOrderHistorys(exchange.FMEX_USDT, 1, 1))
}
