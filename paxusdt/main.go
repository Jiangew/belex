package main

import (
	"fmt"
	"github.com/jiangew/belex/builder"
	"github.com/jiangew/belex/exchange"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	//bot, err := tgbotapi.NewBotAPI("960133387:AAGZ3dZ1FPO-lVJmVTUYsMxDZFUR5WDEEc0")
	//if err != nil {
	//	log.Panic(err)
	//}

	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)

	orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
	if len(orders) > 0 {
		for _, order := range orders {
			cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
			log.Println("cancel order:", order.ID, "ret:", cancel)
		}
	}

	for {
		buyPrice := float64(0)
		sellPrice := float64(0)
		taker, err := api.GetTicker(exchange.PAX_USDT)
		if err != nil {
			log.Println("usdt account got error:", err)
			continue
		} else {
			buyPrice = taker.Buy
			sellPrice = taker.Sell
			log.Println("bid price:", buyPrice)
			log.Println("ask price:", sellPrice)
		}

		orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
		if len(orders) > 0 {
			for _, order := range orders {
				if order.Side == "buy" {
					if order.Price != buyPrice {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("cancel buy:", order.ID, "ret:", cancel)
					}
				} else if order.Side == "sell" {
					if order.Price != sellPrice {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("cancel sell:", order.ID, "ret:", cancel)
					}
				}
			}
		}

		usdtAccount, err := api.GetSubAccount(exchange.USDT)
		if err != nil {
			log.Println("usdt account got error:", err)
		} else {
			if usdtAccount.Available > 10 {
				amount := (usdtAccount.Available - 1) / buyPrice
				if amount > 1 {
					buyOrder, err := api.LimitBuy(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", buyPrice), exchange.PAX_USDT)
					if err != nil {
						log.Println("limit buy amount:", amount, "price:", buyPrice, "error:", err)
					} else {
						log.Println("limit buy amount:", amount, "price:", buyPrice, "success:", buyOrder.ID)
					}
				}
			}
		}

		paxAccount, err := api.GetSubAccount(exchange.PAX)
		if err != nil {
			log.Println("pax account got error:", err)
		} else {
			if paxAccount.Available > 10 {
				amount := paxAccount.Available - 1
				if amount > 1 {
					sellOrder, err := api.LimitSell(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", sellPrice), exchange.PAX_USDT)
					if err != nil {
						log.Println("limit sell amount:", amount, "price:", sellPrice, "error:", err)
					} else {
						log.Println("limit sell amount:", amount, "price:", sellPrice, "success:", sellOrder.ID)
					}
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

//func printfTicker(ticker *exchange.Ticker) {
//	fmt.Println("ticker ", ticker)
//}
//
//func printfDepth(depth *exchange.Depth) {
//	fmt.Println("depth ", depth)
//}
//
//func printfTrade(trade *exchange.Trade) {
//	fmt.Println("trade ", trade)
//}
//
//func printfKline(kline *exchange.Kline, period int) {
//	fmt.Println("kline ", kline)
//}

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
