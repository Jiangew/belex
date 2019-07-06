package main

import (
	"github.com/jiangew/belex/builder"
	"github.com/jiangew/belex/exchange"
	"log"
	"time"
)

func main() {
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)

	log.Println(api.GetTicker(exchange.FMEX_USDT))
	log.Println(api.GetDepth(3, exchange.FMEX_USDT))

	log.Println(api.GetSubAccount(exchange.USDT))
	log.Println(api.GetSubAccount(exchange.FT))
	log.Println(api.GetSubAccount(exchange.FMEX))

	log.Println(api.GetUnfinishOrders(exchange.FMEX_USDT))
}
