### BELEX
BELEX 项目是为了统一并标准化各个数字资产交易平台的接口而设计，同一个策略可以随时切换到任意一个交易平台，而不需要更改任何代码。

### API 接入
```golang
   package main
   
   import (
   	"github.com/jiangew/belex"
   	"github.com/jiangew/belex/builder"
   	"log"
   	"time"
   )
   
   func main() {
   	apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second)
   	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1087")
   	
   	//build spot api
   	//api := apiBuilder.APIKey("").APISecretkey("").ClientID("123").Build(belex.BITSTAMP)
   	api := apiBuilder.APIKey("").APISecretkey("").Build(belex.HUOBI_PRO)
   	log.Println(api.GetExchangeName())
   	log.Println(api.GetTicker(belex.BTC_USD))
   	log.Println(api.GetDepth(2, belex.BTC_USD))
   	//log.Println(api.GetAccount())
   	//log.Println(api.GetUnfinishOrders(belex.BTC_USD))
   
   	//build future api
   	futureApi := apiBuilder.APIKey("").APISecretkey("").BuildFuture(belex.HBDM)
   	log.Println(futureApi.GetExchangeName())
   	log.Println(futureApi.GetFutureTicker(belex.BTC_USD, belex.QUARTER_CONTRACT))
   	log.Println(futureApi.GetFutureDepth(belex.BTC_USD, belex.QUARTER_CONTRACT, 5))
   	//log.Println(futureApi.GetFutureUserinfo()) //account
   	//log.Println(futureApi.GetFuturePosition(belex.BTC_USD , belex.QUARTER_CONTRACT)) //position info
   }
```

### WebSocket 接入
```golang
import (
	"github.com/jiangew/belex"
	"github.com/jiangew/belex/huobi"
	"log"
)

func main() {
	ws := huobi.NewHbdmWs() //huobi期货
	//设置回调函数
	ws.SetCallbacks(func(ticker *belex.FutureTicker) {
		log.Println(ticker)
	}, func(depth *belex.Depth) {
		log.Println(depth)
	}, func(trade *belex.Trade, contract string) {
		log.Println(contract, trade)
	})
	//订阅行情
	ws.SubscribeTrade(belex.BTC_USDT, belex.NEXT_WEEK_CONTRACT)
	ws.SubscribeDepth(belex.BTC_USDT, belex.QUARTER_CONTRACT, 5)
	ws.SubscribeTicker(belex.BTC_USDT, belex.QUARTER_CONTRACT)
}  
```
