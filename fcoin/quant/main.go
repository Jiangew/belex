package main

import (
	"encoding/json"
	"fmt"
	"github.com/jiangew/belex/builder"
	"github.com/jiangew/belex/exchange"
	"time"
)

func main() {
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)

	usdtAccount, _ := api.GetSubAccount(exchange.USDT)
	ret, _ := json.Marshal(usdtAccount)
	fmt.Println(string(ret))

	ftAccount, _ := api.GetSubAccount(exchange.FT)
	ret, _ = json.Marshal(ftAccount)
	fmt.Println(string(ret))

	fmexAccount, _ := api.GetSubAccount(exchange.FMEX)
	ret, _ = json.Marshal(fmexAccount)
	fmt.Println(string(ret))

	ftTicker, _ := api.GetTicker(exchange.FT_USDT)
	ret, _ = json.Marshal(ftTicker)
	fmt.Println(string(ret))

	ftDepth, _ := api.GetDepth(2, exchange.FT_USDT)
	ret, _ = json.Marshal(ftDepth)
	fmt.Println(string(ret))

	fmexTicker, _ := api.GetTicker(exchange.FMEX_USDT)
	ret, _ = json.Marshal(fmexTicker)
	fmt.Println(string(ret))

	fmexDepth, _ := api.GetDepth(2, exchange.FMEX_USDT)
	ret, _ = json.Marshal(fmexDepth)
	fmt.Println(string(ret))

	ftActiveOrders, _ := api.GetActiveOrders(exchange.FT_USDT)
	ret, _ = json.Marshal(ftActiveOrders)
	fmt.Println(string(ret))

	fmexActiveOrders, _ := api.GetActiveOrders(exchange.FMEX_USDT)
	ret, _ = json.Marshal(fmexActiveOrders)
	fmt.Println(string(ret))
}
