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
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)
	buyPrice := float64(0)
	sellPrice := float64(0)
	count := 0

	for {
		if count < 1 {
			orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
			if len(orders) > 0 {
				for i := 0; i < len(orders); {
					order := orders[i]
					cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
					log.Println("id:", order.ID, "cancel:", cancel)
				}
			}
		}

		depths, _ := api.GetDepth(2, exchange.PAX_USDT)
		buyDepth := depths.BidList[0]
		sellDepth := depths.AskList[0]
		log.Println("depth buy:", buyDepth)
		log.Println("depth sell:", sellDepth)

		usdtAccount, err := api.GetSubAccount(exchange.USDT)
		if err != nil {
			log.Println("usdt account err:", err)
		}
		if buyDepth.Price-buyPrice != 0 {
			buyOrders, _ := api.GetActiveOrders(exchange.PAX_USDT)
			if len(buyOrders) > 0 {
				for i := 0; i < len(buyOrders); {
					order := buyOrders[i]
					if (order.Side == "buy" && order.Price == buyPrice) {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("buyid:", order.ID, "cancel:", cancel)
					}
				}
			}
			if usdtAccount.Available > 0 {
				amount := math.Floor(usdtAccount.Available / buyDepth.Price)
				if amount > 1 {
					buyOrder, err := api.LimitBuy(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", buyDepth.Price), exchange.PAX_USDT)
					if err != nil {
						log.Println("limit buy err:", err)
					} else {
						buyPrice = buyDepth.Price
						log.Println("order buy:", buyOrder.ID)
					}
				}
			}
		}

		paxAccount, err := api.GetSubAccount(exchange.PAX)
		if err != nil {
			log.Println("pax account err:", err)
		}
		if sellDepth.Price-sellPrice != 0 {
			sellOrders, _ := api.GetActiveOrders(exchange.PAX_USDT)
			if len(sellOrders) > 0 {
				for i := 0; i < len(sellOrders); {
					order := sellOrders[i]
					if (order.Side == "sell" && order.Price == sellPrice) {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("sellid:", order.ID, "cancel:", cancel)
					}
				}
			}
			if paxAccount.Available > 0 {
				amount := math.Floor(paxAccount.Available / sellDepth.Price)
				if amount > 1 {
					sellOrder, err := api.LimitSell(fmt.Sprintf("%.4f", ), fmt.Sprintf("%.4f", sellDepth.Price), exchange.PAX_USDT)
					if err != nil {
						log.Println("limit sell err:", err)
					} else {
						sellPrice = sellDepth.Price
						log.Println("order sell:", sellOrder.ID)
					}
				}
			}
		}

		count++
		time.Sleep(1 * time.Second)
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
