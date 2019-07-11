package fcoin

import (
	"fmt"
	"github.com/jiangew/belex/exchange"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fcws = NewFCoinWs(&http.Client{
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

func init() {
	fcws.ProxyUrl("socks5://127.0.0.1:1086")
	fcws.SetCallbacks(printfTicker, printfDepth, printfTrade, printfKline)
}

func printfTicker(ticker *exchange.Ticker) {
	fmt.Println(ticker)
}

func printfDepth(depth *exchange.Depth) {
	fmt.Println(depth)
}

func printfTrade(trade *exchange.Trade) {
	fmt.Println(trade)
}

func printfKline(kline *exchange.Kline, period int) {
	fmt.Println(kline)
}

func TestFCoinWs(t *testing.T) {
	fcws.SubscribeTicker(exchange.FMEX_USDT)
	fcws.SubscribeDepth(exchange.FMEX_USDT, 2)
	fcws.SubscribeKline(exchange.FMEX_USDT, exchange.KLINE_PERIOD_1MIN)
	fcws.SubscribeTrade(exchange.FMEX_USDT)
}
