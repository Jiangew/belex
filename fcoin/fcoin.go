package fcoin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jiangew/belex/exchange"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DEPTH_API                 = "market/depth/%s/%s"
	TRADE_URL                 = "orders"
	GET_ACCOUNT_API           = "accounts/balance"
	GET_ORDER_API             = "orders/%s"
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
	SERVER_TIME               = "public/server-time"
)

type FCoinTicker struct {
	exchange.Ticker
	SellAmount,
	BuyAmount float64
}

type FCoin struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey string
	timeoffset   int64
	tradeSymbols []TradeSymbol
}

type TradeSymbol struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
	Tradable      bool   `json:"tradable"`
}

func (fc *FCoin) doAuthenticatedRequest(method, uri string, params url.Values) (interface{}, error) {
	timestamp := time.Now().Unix()*1000 + fc.timeoffset*1000
	sign := fc.buildSigned(method, fc.baseUrl+uri, timestamp, params)

	header := map[string]string{
		"FC-ACCESS-KEY":       fc.accessKey,
		"FC-ACCESS-SIGNATURE": sign,
		"FC-ACCESS-TIMESTAMP": fmt.Sprint(timestamp)}

	var (
		respmap map[string]interface{}
		err     error
	)

	switch method {
	case "GET":
		respmap, err = exchange.HttpGet2(fc.httpClient, fc.baseUrl+uri+"?"+params.Encode(), header)
		if err != nil {
			return nil, err
		}

	case "POST":
		var parammap = make(map[string]string, 1)
		for k, v := range params {
			parammap[k] = v[0]
		}

		respbody, err := exchange.HttpPost(fc.httpClient, fc.baseUrl+uri, parammap, header)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(respbody, &respmap)
	}

	if exchange.ToInt(respmap["status"]) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	return respmap["data"], err
}

func (fc *FCoin) buildSigned(httpmethod string, apiurl string, timestamp int64, para url.Values) string {
	var (
		param = ""
		err   error
	)
	if para != nil {
		param = para.Encode()
	}

	if "GET" == httpmethod && param != "" {
		apiurl += "?" + param
	}

	signStr := httpmethod + apiurl + fmt.Sprint(timestamp)
	if "POST" == httpmethod && param != "" {
		signStr += param
	}

	signStr2, err := url.QueryUnescape(signStr)
	if err != nil {
		signStr2 = signStr
	}

	sign := base64.StdEncoding.EncodeToString([]byte(signStr2))
	mac := hmac.New(sha1.New, []byte(fc.secretKey))
	mac.Write([]byte(sign))
	sum := mac.Sum(nil)
	s := base64.StdEncoding.EncodeToString(sum)

	return s
}

func NewFCoin(client *http.Client, apikey, secretkey string) *FCoin {
	fc := &FCoin{baseUrl: "https://api.fcoin.com/v2/", accessKey: apikey, secretKey: secretkey, httpClient: client}
	_ = fc.setTimeOffset()
	fc.tradeSymbols, _ = fc.getTradeSymbols()

	return fc
}

func (fc *FCoin) GetExchangeName() string {
	return exchange.FCOIN
}

func (fc *FCoin) setTimeOffset() error {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+"public/server-time")
	if err != nil {
		return err
	}

	stime := int64(exchange.ToInt(respmap["data"]))
	st := time.Unix(stime/1000, 0)
	lt := time.Now()
	offset := st.Sub(lt).Seconds()
	fc.timeoffset = int64(offset)

	return nil
}

func (fc *FCoin) GetTicker(symbol exchange.Symbol) (*exchange.Ticker, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/ticker/%s", strings.ToLower(symbol.ToSymbol(""))))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	tick, ok := respmap["data"].(map[string]interface{})
	if !ok {
		return nil, exchange.API_ERR
	}

	tickmap, ok := tick["ticker"].([]interface{})
	if !ok {
		return nil, exchange.API_ERR
	}

	ticker := new(exchange.Ticker)
	ticker.Symbol = symbol.ToSymbol("/")
	ticker.Date = uint64(time.Now().UnixNano() / 1000000)
	ticker.Last = exchange.ToFloat64(tickmap[0])
	ticker.LastVol = exchange.ToFloat64(tickmap[1])
	ticker.Buy = exchange.ToFloat64(tickmap[2])
	ticker.BuyVol = exchange.ToFloat64(tickmap[3])
	ticker.Sell = exchange.ToFloat64(tickmap[4])
	ticker.SellVol = exchange.ToFloat64(tickmap[5])
	ticker.High = exchange.ToFloat64(tickmap[7])
	ticker.Low = exchange.ToFloat64(tickmap[8])
	ticker.Vol = exchange.ToFloat64(tickmap[9])

	return ticker, nil
}

func (fc *FCoin) GetDepth(size int, symbol exchange.Symbol) (*exchange.Depth, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/depth/L20/%s", strings.ToLower(symbol.ToSymbol(""))))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	bids, ok1 := datamap["bids"].([]interface{})
	asks, ok2 := datamap["asks"].([]interface{})
	if !ok1 || !ok2 {
		return nil, errors.New("depth error")
	}

	depth := new(exchange.Depth)
	depth.Symbol = symbol.ToSymbol("/")
	depth.UTime = time.Now()

	n := 0
	for i := 0; i < len(bids); {
		depth.BidList = append(depth.BidList, exchange.DepthRecord{exchange.ToFloat64(bids[i]), exchange.ToFloat64(bids[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	n = 0
	for i := 0; i < len(asks); {
		depth.AskList = append(depth.AskList, exchange.DepthRecord{exchange.ToFloat64(asks[i]), exchange.ToFloat64(asks[i+1])})
		i += 2
		n++
		if n == size {
			break
		}
	}

	//sort.Sort(sort.Reverse(depth.AskList))
	return depth, nil
}

func (fc *FCoin) placeOrder(orderType, orderSide, amount, price string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	params := url.Values{}
	params.Set("side", orderSide)
	params.Set("amount", amount)
	params.Set("symbol", strings.ToLower(symbol.ToSymbol("")))

	switch orderType {
	case "LIMIT", "limit":
		params.Set("type", "limit")
		params.Set("price", price)
	case "MARKET", "market":
		params.Set("type", "market")
	}

	r, err := fc.doAuthenticatedRequest("POST", "orders", params)
	if err != nil {
		return nil, err
	}

	return &exchange.NewOrder{
		ID:        r.(string),
		Symbol:    symbol.ToSymbol("/"),
		Side:      orderSide,
		OrderType: orderType,
		Price:     exchange.ToFloat64(price),
		Amount:    exchange.ToFloat64(amount),
		State:     "SUBMITTED",
	}, nil
}

func (fc *FCoin) LimitBuy(amount, price string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	return fc.placeOrder("limit", "buy", amount, price, symbol)
}

func (fc *FCoin) LimitSell(amount, price string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	return fc.placeOrder("limit", "sell", amount, price, symbol)
}

func (fc *FCoin) MarketBuy(amount, price string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	return fc.placeOrder("market", "buy", amount, price, symbol)
}

func (fc *FCoin) MarketSell(amount, price string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	return fc.placeOrder("market", "sell", amount, price, symbol)
}

func (fc *FCoin) CancelOrder(orderId string, symbol exchange.Symbol) (bool, error) {
	uri := fmt.Sprintf("orders/%s/submit-cancel", orderId)
	_, err := fc.doAuthenticatedRequest("POST", uri, url.Values{})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (fc *FCoin) toOrder(o map[string]interface{}, symbol exchange.Symbol) *exchange.NewOrder {
	var fees float64
	refund := exchange.ToFloat64(o["fees_income"])
	fee := exchange.ToFloat64(o["fill_fees"])
	if fee == 0 {
		fees = -refund
	} else {
		fees = fee
	}

	return &exchange.NewOrder{
		ID:           o["id"].(string),
		Symbol:       symbol.ToSymbol("/"),
		Side:         o["side"].(string),
		OrderType:    o["type"].(string),
		Price:        exchange.ToFloat64(o["price"]),
		Amount:       exchange.ToFloat64(o["amount"]),
		FilledAmount: exchange.ToFloat64(o["filled_amount"]),
		FillFee:      fees,
		State:        o["state"].(string),
		CreatedAt:    exchange.ToUint64(o["created_at"]),
	}
}

func (fc *FCoin) GetOrder(orderId string, symbol exchange.Symbol) (*exchange.NewOrder, error) {
	uri := fmt.Sprintf("orders/%s", orderId)
	r, err := fc.doAuthenticatedRequest("GET", uri, url.Values{})

	if err != nil {
		return nil, err
	}

	return fc.toOrder(r.(map[string]interface{}), symbol), nil
}

func (fc *FCoin) GetActiveOrders(symbol exchange.Symbol) ([]exchange.NewOrder, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(symbol.ToSymbol("")))
	params.Set("states", "submitted,partial_filled")
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.NewOrder
	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), symbol))
	}

	return ords, nil
}

func (fc *FCoin) getAfterTimeOrderHistorys(symbol exchange.Symbol, times time.Time) ([]exchange.NewOrder, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(symbol.ToSymbol("")))
	params.Set("states", "filled")
	params.Set("after", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.NewOrder

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), symbol))
	}

	return ords, nil
}

func (fc *FCoin) getBeforeTimeOrderHistorys(symbol exchange.Symbol, times time.Time) ([]exchange.NewOrder, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(symbol.ToSymbol("")))
	params.Set("states", "filled")
	params.Set("before", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.NewOrder

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), symbol))
	}

	return ords, nil
}

func (fc *FCoin) GetHoursOrderHistorys(symbol exchange.Symbol, start time.Time, hours int64) ([]exchange.NewOrder, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(symbol, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(symbol, start.Add(time.Hour*time.Duration(hours)))

	ords := make([]exchange.NewOrder, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.ID == v2.ID {
				ords = append(ords, v1)
			}
		}
	}

	return ords, nil
}

func (fc *FCoin) GetDaysOrderHistorys(symbol exchange.Symbol, start time.Time, days int64) ([]exchange.NewOrder, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(symbol, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(symbol, start.Add(time.Hour*24*time.Duration(days)))

	ords := make([]exchange.NewOrder, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.ID == v2.ID {
				ords = append(ords, v1)
			}
		}
	}

	if len(ords) == 0 && len(ord2) > 1 && len(ord1) > 1 {
		return nil, errors.New("more than 100 orders")
	}

	return ords, nil
}

func (fc *FCoin) GetOrderHistorys(symbol exchange.Symbol, currentPage, pageSize int) ([]exchange.NewOrder, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(symbol.ToSymbol("")))
	params.Set("states", "partial_canceled,filled")
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.NewOrder
	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), symbol))
	}

	return ords, nil
}

func (fc *FCoin) GetAccount() (*exchange.Account, error) {
	r, err := fc.doAuthenticatedRequest("GET", "accounts/balance", url.Values{})
	if err != nil {
		return nil, err
	}

	acc := new(exchange.Account)
	acc.SubAccounts = make(map[exchange.Currency]exchange.SubAccount)

	balances := r.([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := exchange.NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = exchange.SubAccount{
			Currency:      currency.Name,
			Available:     exchange.ToFloat64(vv["available"]),
			Frozen:        exchange.ToFloat64(vv["frozen"]),
			DemandDeposit: exchange.ToFloat64(vv["demand_deposit"]),
			LockDeposit:   exchange.ToFloat64(vv["lock_deposit"]),
			Balance:       exchange.ToFloat64(vv["balance"]),
		}
	}

	return acc, nil
}

func (fc *FCoin) GetSubAccount(currency exchange.Currency) (*exchange.SubAccount, error) {
	r, err := fc.doAuthenticatedRequest("GET", "accounts/balance", url.Values{})
	if err != nil {
		return nil, err
	}

	subaccount := new(exchange.SubAccount)
	balances := r.([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		if strings.ToLower(currency.Name) != vv["currency"].(string) {
			continue
		}
		subaccount.Currency = currency.Name
		subaccount.Available = exchange.ToFloat64(vv["available"])
		subaccount.Frozen = exchange.ToFloat64(vv["frozen"])
		subaccount.DemandDeposit = exchange.ToFloat64(vv["demand_deposit"])
		subaccount.LockDeposit = exchange.ToFloat64(vv["lock_deposit"])
		subaccount.Balance = exchange.ToFloat64(vv["balance"])
	}

	return subaccount, nil
}

func (fc *FCoin) GetKlines(symbol exchange.Symbol) ([]exchange.Kline, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/candles/M1/%s?limit=10&before=%d", strings.ToLower(symbol.ToSymbol("")), time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].([]interface{})
	klines := make([]exchange.Kline, 0)
	for _, v := range datamap {
		vv := v.(map[string]interface{})
		var kline exchange.Kline
		kline.Symbol = symbol
		kline.Timestamp = int64(vv["id"].(float64))
		kline.Vol = vv["quote_vol"].(float64)
		kline.Open = vv["open"].(float64)
		kline.High = vv["high"].(float64)
		kline.Low = vv["low"].(float64)
		kline.Close = vv["close"].(float64)
		klines = append(klines, kline)
	}

	return klines, nil
}

func (fc *FCoin) IsOrderable(symbol exchange.Symbol) (bool, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/candles/M5/%s?limit=6&before=%d", strings.ToLower(symbol.ToSymbol("")), time.Now().Unix()))
	if err != nil {
		return false, err
	}

	if respmap["status"].(float64) != 0 {
		return false, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].([]interface{})
	var high string
	var low string
	isOrderable := false
	highEq := false
	lowEq := false
	curse := 1

	for _, v := range datamap {
		vv := v.(map[string]interface{})
		//id := int64(vv["id"].(float64))
		curHigh := exchange.FloatToString(vv["high"].(float64), 4)
		curLow := exchange.FloatToString(vv["low"].(float64), 4)

		if curse == 1 {
			high = curHigh
			low = curLow
		} else {
			if high == curHigh {
				highEq = true
			}
			if low == curLow {
				lowEq = true
			}
			//log.Println(id, high, curHigh, low, curLow)
		}

		curse = curse + 1
	}

	if highEq && lowEq {
		isOrderable = true
	}

	return isOrderable, nil
}

func (fc *FCoin) GetTrades(symbol exchange.Symbol, since int64) ([]exchange.Trade, error) {
	panic("todo implement")
}

func (fc *FCoin) getTradeSymbols() ([]TradeSymbol, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+"public/symbols")
	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	datamap := respmap["data"].([]interface{})
	tradeSymbols := make([]TradeSymbol, 0)
	for _, v := range datamap {
		vv := v.(map[string]interface{})
		var symbol TradeSymbol
		symbol.Name = vv["name"].(string)
		symbol.BaseCurrency = vv["base_currency"].(string)
		symbol.QuoteCurrency = vv["quote_currency"].(string)
		symbol.PriceDecimal = int(vv["price_decimal"].(float64))
		symbol.AmountDecimal = int(vv["amount_decimal"].(float64))
		symbol.Tradable = vv["tradable"].(bool)

		if symbol.Tradable {
			tradeSymbols = append(tradeSymbols, symbol)
		}
	}

	return tradeSymbols, nil
}

func (fc *FCoin) GetTradeSymbols(symbol exchange.Symbol) (*TradeSymbol, error) {
	if len(fc.tradeSymbols) == 0 {
		var err error
		fc.tradeSymbols, err = fc.getTradeSymbols()
		if err != nil {
			return nil, err
		}
	}

	for k, v := range fc.tradeSymbols {
		if v.Name == strings.ToLower(symbol.ToSymbol("")) {
			return &fc.tradeSymbols[k], nil
		}
	}

	return nil, errors.New("symbol not found")
}
