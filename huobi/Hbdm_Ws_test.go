package huobi

import (
	"github.com/jiangew/belex"
	"log"
	"testing"
	"time"
)

func TestNewHbdmWs(t *testing.T) {
	ws := NewHbdmWs()
	ws.ProxyUrl("socks5://127.0.0.1:1080")

	ws.SetCallbacks(func(ticker *belex.FutureTicker) {
		log.Println(ticker.Ticker)
	}, func(depth *belex.Depth) {
		log.Println(">>>>>>>>>>>>>>>")
		log.Println(depth.ContractType, depth.Pair)
		log.Println(depth.BidList)
		log.Println(depth.AskList)
		log.Println("<<<<<<<<<<<<<<")
	}, func(trade *belex.Trade, s string) {
		log.Println(s, trade)
	})

	t.Log(ws.SubscribeTicker(belex.BTC_USD, belex.QUARTER_CONTRACT))
	t.Log(ws.SubscribeDepth(belex.BTC_USD, belex.NEXT_WEEK_CONTRACT, 0))
	t.Log(ws.SubscribeTrade(belex.LTC_USD, belex.THIS_WEEK_CONTRACT))
	time.Sleep(time.Minute)
}
