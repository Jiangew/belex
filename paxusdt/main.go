package main

import (
	"fmt"
	"github.com/jiangew/belex/builder"
	"github.com/jiangew/belex/exchange"
	"log"
	"math"
	"time"
)

func printfTicker(ticker *exchange.Ticker) {
	fmt.Println("ticker ", ticker)
}

func printfDepth(depth *exchange.Depth) {
	fmt.Println("depth ", depth)
}

func printfTrade(trade *exchange.Trade) {
	fmt.Println("trade ", trade)
}

func printfKline(kline *exchange.Kline, period int) {
	fmt.Println("kline ", kline)
}

func main() {
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)
	buyPrice := float64(0)
	sellPrice := float64(0)
	buyID := ""
	sellID := ""
	count := 0

	for {
		if count < 1 {
			paxOrders, _ := api.GetActiveOrders(exchange.PAX_USDT)
			if len(paxOrders) > 0 {
				for i := 0; i < len(paxOrders); {
					order := paxOrders[i]
					cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
					log.Println("id:", order.ID, "cancel:", cancel)
				}
			}
		}

		usdtAccount, _ := api.GetSubAccount(exchange.USDT)
		paxAccount, _ := api.GetSubAccount(exchange.PAX)

		//paxTicker, _ := api.GetTicker(exchange.PAX_USDT)
		//ret, _ = json.Marshal(paxTicker)
		//log.Println(string(ret))

		depths, _ := api.GetDepth(2, exchange.PAX_USDT)
		buyDepth := depths.BidList[0]
		sellDepth := depths.AskList[0]
		log.Println("buyDepth:", buyDepth)
		log.Println("sellDepth:", sellDepth)

		if buyDepth.Price-buyPrice != 0 {
			if buyID != "" {
				cancel, _ := api.CancelOrder(buyID, exchange.PAX_USDT)
				log.Println("id:", buyID, "cancel:", cancel)
			}

			if usdtAccount.Available > 0 {
				buyOrder, err := api.LimitBuy(fmt.Sprintf("%.4f", math.Floor(usdtAccount.Available/buyDepth.Price)), fmt.Sprintf("%.4f", buyDepth.Price), exchange.PAX_USDT)
				if err != nil {
					log.Println(err)
				} else {
					buyPrice = buyDepth.Price
					buyID = buyOrder.ID
					log.Println("buyOrderID:", buyOrder.ID)
				}
			}
		}

		if sellDepth.Price-sellPrice != 0 {
			if sellID != "" {
				cancel, _ := api.CancelOrder(sellID, exchange.PAX_USDT)
				log.Println("id:", sellID, "cancel:", cancel)
			}

			if paxAccount.Available > 0 {
				sellOrder, err := api.LimitSell(fmt.Sprintf("%.4f", math.Floor(paxAccount.Available/sellDepth.Price)), fmt.Sprintf("%.4f", sellDepth.Price), exchange.PAX_USDT)
				if err != nil {
					log.Println(err)
				} else {
					sellPrice = sellDepth.Price
					sellID = sellOrder.ID
					log.Println("sellOrderID:", sellOrder.ID)
				}
			}
		}

		count++
	}
}

//func main() {
//	fcws := fcoin.NewFCoinWs(&http.Client{
//		Transport: &http.Transport{
//			Proxy: func(req *http.Request) (*url.URL, error) {
//				return url.Parse("socks5://127.0.0.1:1086")
//				return nil, nil
//			},
//			Dial: (&net.Dialer{
//				Timeout: 10 * time.Second,
//			}).Dial,
//		},
//		Timeout: 10 * time.Second,
//	}, "1412ac27e3f741c796f7c4600069d9f1", "4843754749be46919d986142917f06d7")
//
//	fcws.ProxyUrl("socks5://127.0.0.1:1086")
//	fcws.SetCallbacks(printfTicker, printfDepth, printfTrade, printfKline)
//
//	fcws.SubscribeTicker(exchange.PAX_USDT)
//	fcws.SubscribeDepth(exchange.PAX_USDT, 2)
//	fcws.SubscribeKline(exchange.PAX_USDT, exchange.KLINE_PERIOD_1MIN)
//	fcws.SubscribeTrade(exchange.PAX_USDT)
//
//	time.Sleep(60 * 60 * time.Second)
//}
