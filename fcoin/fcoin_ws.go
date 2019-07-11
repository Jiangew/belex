package fcoin

import (
	"errors"
	"fmt"
	"github.com/jiangew/belex/exchange"
	"github.com/json-iterator/go"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	FCoinWSTicker        = "ticker.%s"
	FCoinWSOrderBook     = "depth.L%d.%s"
	FCoinWSOrderBookL20  = "depth.L20.%s"
	FCoinWSOrderBookL150 = "depth.L150.%s"
	FCoinWSOrderBookFull = "depth.full.%s"
	FCoinWSTrades        = "trade.%s"
	FCoinWSKLines        = "candle.%s.%s"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FCoinWs struct {
	*exchange.WsBuilder
	sync.Once
	wsConn *exchange.WsConn

	tickerCallback func(*exchange.Ticker)
	depthCallback  func(*exchange.Depth)
	tradeCallback  func(*exchange.Trade)
	klineCallback  func(*exchange.Kline, int)

	clientId      string
	subcribeTypes []string
	timeoffset    int64
	tradeSymbols  []TradeSymbol
}

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	exchange.KLINE_PERIOD_1MIN:   "M1",
	exchange.KLINE_PERIOD_3MIN:   "M3",
	exchange.KLINE_PERIOD_5MIN:   "M5",
	exchange.KLINE_PERIOD_15MIN:  "M15",
	exchange.KLINE_PERIOD_30MIN:  "M30",
	exchange.KLINE_PERIOD_60MIN:  "H1",
	exchange.KLINE_PERIOD_4H:     "H4",
	exchange.KLINE_PERIOD_6H:     "H6",
	exchange.KLINE_PERIOD_1DAY:   "D1",
	exchange.KLINE_PERIOD_1WEEK:  "W1",
	exchange.KLINE_PERIOD_1MONTH: "MN",
}
var _INERNAL_KLINE_PERIOD_REVERTER = map[string]int{
	"M1":  exchange.KLINE_PERIOD_1MIN,
	"M3":  exchange.KLINE_PERIOD_3MIN,
	"M5":  exchange.KLINE_PERIOD_5MIN,
	"M15": exchange.KLINE_PERIOD_15MIN,
	"M30": exchange.KLINE_PERIOD_30MIN,
	"H1":  exchange.KLINE_PERIOD_60MIN,
	"H4":  exchange.KLINE_PERIOD_4H,
	"H6":  exchange.KLINE_PERIOD_6H,
	"D1":  exchange.KLINE_PERIOD_1DAY,
	"W1":  exchange.KLINE_PERIOD_1WEEK,
	"MN":  exchange.KLINE_PERIOD_1MONTH,
}

func NewFCoinWs(client *http.Client, apikey, secretkey string) *FCoinWs {
	fcWs := &FCoinWs{}
	fcWs.clientId = getRandomString(8)
	fcWs.WsBuilder = exchange.NewWsBuilder().
		WsUrl("wss://api.fcoin.com/v2/ws").
		Heartbeat2(func() interface{} {
			ts := time.Now().Unix()*1000 + fcWs.timeoffset*1000
			args := make([]interface{}, 0)
			args = append(args, ts)
			return map[string]interface{}{
				"cmd":  "ping",
				"id":   fcWs.clientId,
				"args": args}

		}, 25*time.Second).
		ReconnectIntervalTime(24 * time.Hour).
		UnCompressFunc(exchange.FlateUnCompress).
		ProtoHandleFunc(fcWs.handle)
	fc := NewFCoin(client, apikey, secretkey)
	fcWs.tradeSymbols = fc.tradeSymbols

	if len(fcWs.tradeSymbols) == 0 {
		panic("trade symbol is empty, pls check connection...")
	}

	return fcWs
}

func getRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

func (fcWs *FCoinWs) SetCallbacks(
	tickerCallback func(*exchange.Ticker),
	depthCallback func(*exchange.Depth),
	tradeCallback func(*exchange.Trade),
	klineCallback func(*exchange.Kline, int),
) {
	fcWs.tickerCallback = tickerCallback
	fcWs.depthCallback = depthCallback
	fcWs.tradeCallback = tradeCallback
	fcWs.klineCallback = klineCallback
}

func (fcWs *FCoinWs) subscribe(sub map[string]interface{}) error {
	fcWs.connectWs()
	return fcWs.wsConn.Subscribe(sub)
}

func (fcWs *FCoinWs) SubscribeDepth(symbol exchange.Symbol, size int) error {
	if fcWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	arg := fmt.Sprintf(FCoinWSOrderBook, size, strings.ToLower(symbol.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeTicker(symbol exchange.Symbol) error {
	if fcWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	arg := fmt.Sprintf(FCoinWSTicker, strings.ToLower(symbol.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeTrade(symbol exchange.Symbol) error {
	if fcWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	arg := fmt.Sprintf(FCoinWSTrades, strings.ToLower(symbol.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) SubscribeKline(symbol exchange.Symbol, period int) error {
	if fcWs.klineCallback == nil {
		return errors.New("place set kline callback func")
	}
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "M1"
	}

	arg := fmt.Sprintf(FCoinWSKLines, periodS, strings.ToLower(symbol.ToSymbol("")))
	args := make([]interface{}, 0)
	args = append(args, arg)

	return fcWs.subscribe(map[string]interface{}{
		"id":   fcWs.clientId,
		"cmd":  "sub",
		"args": args})
}

func (fcWs *FCoinWs) connectWs() {
	fcWs.Do(func() {
		fcWs.wsConn = fcWs.WsBuilder.Build()
		fcWs.wsConn.ReceiveMessage()
	})
}

func (fcWs *FCoinWs) parseTickerData(tickmap []interface{}) *exchange.Ticker {
	t := new(exchange.Ticker)
	t.Date = uint64(time.Now().UnixNano() / 1000000)
	t.Last = exchange.ToFloat64(tickmap[0])
	t.Vol = exchange.ToFloat64(tickmap[9])
	t.Low = exchange.ToFloat64(tickmap[8])
	t.High = exchange.ToFloat64(tickmap[7])
	t.Buy = exchange.ToFloat64(tickmap[2])
	t.Sell = exchange.ToFloat64(tickmap[4])

	return t
}

func (fcWs *FCoinWs) parseDepthData(bids, asks []interface{}) *exchange.Depth {
	depth := new(exchange.Depth)
	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, exchange.DepthRecord{exchange.ToFloat64(bids[i]), exchange.ToFloat64(bids[i+1])})
		i += 2
		n++
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, exchange.DepthRecord{exchange.ToFloat64(asks[i]), exchange.ToFloat64(asks[i+1])})
		i += 2
		n++
	}

	return depth
}

func (fcWs *FCoinWs) parseKlineData(tickmap []interface{}) *exchange.Ticker {
	t := new(exchange.Ticker)
	t.Date = uint64(time.Now().UnixNano() / 1000000)
	t.Last = exchange.ToFloat64(tickmap[0])
	t.Vol = exchange.ToFloat64(tickmap[9])
	t.Low = exchange.ToFloat64(tickmap[8])
	t.High = exchange.ToFloat64(tickmap[7])
	t.Buy = exchange.ToFloat64(tickmap[2])
	t.Sell = exchange.ToFloat64(tickmap[4])

	return t
}

func (fcWs *FCoinWs) handle(msg []byte) error {
	datamap := make(map[string]interface{})
	err := json.Unmarshal(msg, &datamap)
	if err != nil {
		fmt.Println("json unmarshal error for ", string(msg))
		return err
	}

	msgType, isOk := datamap["type"].(string)
	if isOk {
		resp := strings.Split(msgType, ".")
		switch resp[0] {
		case "hello", "ping":
			fcWs.wsConn.UpdateActiveTime()
			stime := int64(exchange.ToInt(datamap["ts"]))
			st := time.Unix(0, stime*1000*1000)
			lt := time.Now()
			offset := st.Sub(lt).Seconds()
			fcWs.timeoffset = int64(offset)
		case "ticker":
			tick := fcWs.parseTickerData(datamap["ticker"].([]interface{}))
			symbol, err := fcWs.getSymbolFromType(resp[1])
			if err != nil {
				panic(err)
			}
			tick.Symbol = symbol.ToSymbol("/")
			fcWs.tickerCallback(tick)
			return nil
		case "depth":
			dep := fcWs.parseDepthData(datamap["bids"].([]interface{}), datamap["asks"].([]interface{}))
			stime := int64(exchange.ToInt(datamap["ts"]))
			dep.UTime = time.Unix(stime/1000, 0)
			symbol, err := fcWs.getSymbolFromType(resp[2])
			if err != nil {
				panic(err)
			}
			dep.Symbol = symbol.ToSymbol("/")

			fcWs.depthCallback(dep)
			return nil
		case "candle":
			period := _INERNAL_KLINE_PERIOD_REVERTER[resp[1]]
			kline := &exchange.Kline{
				Timestamp: int64(exchange.ToInt(datamap["id"])),
				Open:      exchange.ToFloat64(datamap["open"]),
				Close:     exchange.ToFloat64(datamap["close"]),
				High:      exchange.ToFloat64(datamap["high"]),
				Low:       exchange.ToFloat64(datamap["low"]),
				Vol:       exchange.ToFloat64(datamap["quote_vol"]),
			}
			symbol, err := fcWs.getSymbolFromType(resp[2])
			if err != nil {
				panic(err)
			}
			kline.Symbol = symbol
			fcWs.klineCallback(kline, period)
			return nil
		case "trade":
			side := exchange.BUY
			if datamap["side"] == "sell" {
				side = exchange.SELL
			}
			trade := &exchange.Trade{
				Tid:    int64(exchange.ToUint64(datamap["id"])),
				Type:   exchange.TradeSide(side),
				Amount: exchange.ToFloat64(datamap["amount"]),
				Price:  exchange.ToFloat64(datamap["price"]),
				Date:   int64(exchange.ToUint64(datamap["ts"])),
			}
			symbol, err := fcWs.getSymbolFromType(resp[1])
			if err != nil {
				panic(err)
			}
			trade.Symbol = symbol
			fcWs.tradeCallback(trade)
			return nil
		default:
			return errors.New("unknown message " + msgType)
		}
	}

	return nil
}

func (fcWs *FCoinWs) getSymbolFromType(symbol string) (exchange.Symbol, error) {
	for _, v := range fcWs.tradeSymbols {
		if v.Name == symbol {
			return exchange.NewSymbol2(v.BaseCurrency + "_" + v.QuoteCurrency), nil
		}
	}

	return exchange.NewSymbol2("" + "_" + ""), errors.New("symbol not support :" + symbol)
}
