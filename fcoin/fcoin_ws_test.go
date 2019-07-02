package fcoin

import (
	"fmt"
	"github.com/jiangew/belex"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var fcws = NewFCoinWs(&http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("socks5://127.0.0.1:1080")
			return nil, nil
		},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	},
	Timeout: 10 * time.Second,
})

func init() {
	fcws.ProxyUrl("socks5://127.0.0.1:1080")
	fcws.SetCallbacks(printfTicker, printfDepth, printfTrade, printfKline)
}

func printfTicker(ticker *belex.Ticker) {
	fmt.Println(ticker)
}
func printfDepth(depth *belex.Depth) {
	fmt.Println(depth)
}
func printfTrade(trade *belex.Trade) {
	fmt.Println(trade)
}
func printfKline(kline *belex.Kline, period int) {
	fmt.Println(kline)
}

func TestFCoinWs_GetTickerWithWs(t *testing.T) {
	return
	fcws.SubscribeTicker(belex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
func TestFCoinWs_GetDepthWithWs(t *testing.T) {
	return
	fcws.SubscribeDepth(belex.BTC_USDT, 20)
	time.Sleep(time.Second * 10)
}
func TestFCoinWs_GetKLineWithWs(t *testing.T) {
	//return
	fcws.SubscribeKline(belex.BTC_USDT, belex.KLINE_PERIOD_1MIN)
	time.Sleep(time.Second * 120)
}
func TestFCoinWs_GetTradesWithWs(t *testing.T) {
	return
	fcws.SubscribeTrade(belex.BTC_USDT)
	time.Sleep(time.Second * 10)
}
