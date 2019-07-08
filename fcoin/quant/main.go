package main

import (
	"encoding/json"
	"github.com/jiangew/belex/builder"
	"github.com/jiangew/belex/exchange"
	"log"
	"time"
)

func main() {
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)

	usdtAccount, _ := api.GetSubAccount(exchange.USDT)
	ret, _ := json.Marshal(usdtAccount)
	log.Println(string(ret))

	tusdAccount, _ := api.GetSubAccount(exchange.TUSD)
	ret, _ = json.Marshal(tusdAccount)
	log.Println(string(ret))

	tusdTicker, _ := api.GetTicker(exchange.TUSD_USDT)
	ret, _ = json.Marshal(tusdTicker)
	log.Println(string(ret))

	tusdDepth, _ := api.GetDepth(2, exchange.TUSD_USDT)
	ret, _ = json.Marshal(tusdDepth)
	log.Println(string(ret))

	tusdActiveOrders, _ := api.GetActiveOrders(exchange.TUSD_USDT)
	ret, _ = json.Marshal(tusdActiveOrders)
	log.Println(string(ret))
}
