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
	DEPTH_API       = "market/depth/%s/%s"
	TRADE_URL       = "orders"
	GET_ACCOUNT_API = "accounts/balance"
	GET_ORDER_API   = "orders/%s"
	//GET_ORDERS_LIST_API             = ""
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

func NewFCoin(client *http.Client, apikey, secretkey string) *FCoin {
	fc := &FCoin{baseUrl: "https://api.fcoin.com/v2/", accessKey: apikey, secretKey: secretkey, httpClient: client}
	fc.setTimeOffset()
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

func (fc *FCoin) GetTicker(currencyPair exchange.CurrencyPair) (*exchange.Ticker, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/ticker/%s",
		strings.ToLower(currencyPair.ToSymbol(""))))

	if err != nil {
		return nil, err
	}

	if respmap["status"].(float64) != 0 {
		return nil, errors.New(respmap["msg"].(string))
	}

	//
	tick, ok := respmap["data"].(map[string]interface{})
	if !ok {
		return nil, exchange.API_ERR
	}

	tickmap, ok := tick["ticker"].([]interface{})
	if !ok {
		return nil, exchange.API_ERR
	}

	ticker := new(exchange.Ticker)
	ticker.Pair = currencyPair
	ticker.Date = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	ticker.Last = exchange.ToFloat64(tickmap[0])
	ticker.Vol = exchange.ToFloat64(tickmap[9])
	ticker.Low = exchange.ToFloat64(tickmap[8])
	ticker.High = exchange.ToFloat64(tickmap[7])
	ticker.Buy = exchange.ToFloat64(tickmap[2])
	ticker.Sell = exchange.ToFloat64(tickmap[4])

	return ticker, nil
}

func (fc *FCoin) GetDepth(size int, currency exchange.CurrencyPair) (*exchange.Depth, error) {
	respmap, err := exchange.HttpGet(fc.httpClient, fc.baseUrl+fmt.Sprintf("market/depth/L20/%s", strings.ToLower(currency.ToSymbol(""))))
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
	depth.Pair = currency

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

		respbody, err := exchange.HttpPostForm4(fc.httpClient, fc.baseUrl+uri, parammap, header)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(respbody, &respmap)
	}

	//log.Println(respmap)
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

	signStr2, err := url.QueryUnescape(signStr) // 不需要编码
	if err != nil {
		signStr2 = signStr
	}

	sign := base64.StdEncoding.EncodeToString([]byte(signStr2))

	mac := hmac.New(sha1.New, []byte(fc.secretKey))

	mac.Write([]byte(sign))
	sum := mac.Sum(nil)

	s := base64.StdEncoding.EncodeToString(sum)
	//log.Println(s)
	return s
}

func (fc *FCoin) placeOrder(orderType, orderSide, amount, price string, pair exchange.CurrencyPair) (*exchange.Order, error) {
	params := url.Values{}

	params.Set("side", orderSide)
	params.Set("amount", amount)
	//params.Set("price", price)
	params.Set("symbol", strings.ToLower(pair.AdaptUsdToUsdt().ToSymbol("")))

	switch orderType {
	case "LIMIT", "limit":
		params.Set("price", price)
		params.Set("type", "limit")
	case "MARKET", "market":
		params.Set("type", "market")
	}

	r, err := fc.doAuthenticatedRequest("POST", "orders", params)
	if err != nil {
		return nil, err
	}

	side := exchange.SELL
	if orderSide == "buy" {
		side = exchange.BUY
	}

	return &exchange.Order{
		Currency: pair,
		OrderID2: r.(string),
		Amount:   exchange.ToFloat64(amount),
		Price:    exchange.ToFloat64(price),
		Side:     exchange.TradeSide(side),
		Status:   exchange.ORDER_UNFINISH}, nil
}

func (fc *FCoin) LimitBuy(amount, price string, currency exchange.CurrencyPair) (*exchange.Order, error) {
	return fc.placeOrder("limit", "buy", amount, price, currency)
}

func (fc *FCoin) LimitSell(amount, price string, currency exchange.CurrencyPair) (*exchange.Order, error) {
	return fc.placeOrder("limit", "sell", amount, price, currency)
}

func (fc *FCoin) MarketBuy(amount, price string, currency exchange.CurrencyPair) (*exchange.Order, error) {
	return fc.placeOrder("market", "buy", amount, price, currency)
}

func (fc *FCoin) MarketSell(amount, price string, currency exchange.CurrencyPair) (*exchange.Order, error) {
	return fc.placeOrder("market", "sell", amount, price, currency)
}

func (fc *FCoin) CancelOrder(orderId string, currency exchange.CurrencyPair) (bool, error) {
	uri := fmt.Sprintf("orders/%s/submit-cancel", orderId)
	_, err := fc.doAuthenticatedRequest("POST", uri, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fc *FCoin) toOrder(o map[string]interface{}, pair exchange.CurrencyPair) *exchange.Order {
	side := exchange.SELL
	if o["side"].(string) == "buy" {
		side = exchange.BUY
	}

	orderStatus := exchange.ORDER_UNFINISH
	switch o["state"].(string) {
	case "partial_filled":
		orderStatus = exchange.ORDER_PART_FINISH
	case "filled":
		orderStatus = exchange.ORDER_FINISH
	case "pending_cancel":
		orderStatus = exchange.ORDER_CANCEL_ING
	case "canceled", "partial_canceled":
		orderStatus = exchange.ORDER_CANCEL
	}
	var fees float64
	refund := exchange.ToFloat64(o["fees_income"])
	fee := exchange.ToFloat64(o["fill_fees"])
	if fee == 0 {
		fees = -refund
	} else {
		fees = fee
	}
	return &exchange.Order{
		Currency:   pair,
		Side:       exchange.TradeSide(side),
		OrderID2:   o["id"].(string),
		Amount:     exchange.ToFloat64(o["amount"]),
		Price:      exchange.ToFloat64(o["price"]),
		DealAmount: exchange.ToFloat64(o["filled_amount"]),
		Status:     exchange.TradeStatus(orderStatus),
		Fee:        fees,
		OrderTime:  exchange.ToInt(o["created_at"])}
}

func (fc *FCoin) GetOneOrder(orderId string, currency exchange.CurrencyPair) (*exchange.Order, error) {
	uri := fmt.Sprintf("orders/%s", orderId)
	r, err := fc.doAuthenticatedRequest("GET", uri, url.Values{})

	if err != nil {
		return nil, err
	}

	return fc.toOrder(r.(map[string]interface{}), currency), nil

}

func (fc *FCoin) GetOrdersList() {
	//path := API_URL + fmt.Sprintf(CANCEL_ORDER_API, strings.ToLower(currency.ToSymbol("")))

}

func (fc *FCoin) GetUnfinishOrders(currency exchange.CurrencyPair) ([]exchange.Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "submitted,partial_filled")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fc *FCoin) getAfterTimeOrderHistorys(currency exchange.CurrencyPair, times time.Time) ([]exchange.Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "filled")
	params.Set("after", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}
func (fc *FCoin) getBeforeTimeOrderHistorys(currency exchange.CurrencyPair, times time.Time) ([]exchange.Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "filled")
	params.Set("before", fmt.Sprint(times.Unix()*1000))
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}

	var ords []exchange.Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
	}

	return ords, nil
}

func (fc *FCoin) GetHoursOrderHistorys(currency exchange.CurrencyPair, start time.Time, hours int64) ([]exchange.Order, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(currency, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(currency, start.Add(time.Hour*time.Duration(hours)))
	ords := make([]exchange.Order, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.OrderID2 == v2.OrderID2 {
				ords = append(ords, v1)
			}
		}
	}
	return ords, nil
}

func (fc *FCoin) GetDaysOrderHistorys(currency exchange.CurrencyPair, start time.Time, days int64) ([]exchange.Order, error) {
	ord1, _ := fc.getAfterTimeOrderHistorys(currency, start)
	ord2, _ := fc.getBeforeTimeOrderHistorys(currency, start.Add(time.Hour*24*time.Duration(days)))
	ords := make([]exchange.Order, 0)
	for _, v1 := range ord1 {
		for _, v2 := range ord2 {
			if v1.OrderID2 == v2.OrderID2 {
				ords = append(ords, v1)
			}
		}
	}
	if len(ords) == 0 && len(ord2) > 1 && len(ord1) > 1 {
		return nil, errors.New("more than 100 orders")
	}
	return ords, nil
}

func (fc *FCoin) GetOrderHistorys(currency exchange.CurrencyPair, currentPage, pageSize int) ([]exchange.Order, error) {
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol("")))
	params.Set("states", "partial_canceled,filled")
	//params.Set("before", "1")
	//params.Set("after", "0")
	params.Set("limit", "100")

	r, err := fc.doAuthenticatedRequest("GET", "orders", params)
	if err != nil {
		return nil, err
	}
	var ords []exchange.Order

	for _, ord := range r.([]interface{}) {
		ords = append(ords, *fc.toOrder(ord.(map[string]interface{}), currency))
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
	acc.Exchange = fc.GetExchangeName()

	balances := r.([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := exchange.NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = exchange.SubAccount{
			Currency:     currency,
			Amount:       exchange.ToFloat64(vv["available"]),
			ForzenAmount: exchange.ToFloat64(vv["frozen"]),
		}
	}

	return acc, nil

}

func (fc *FCoin) GetKlineRecords(currency exchange.CurrencyPair, period, size, since int) ([]exchange.Kline, error) {
	panic("not implement")
}

//非个人，整个交易所的交易记录
func (fc *FCoin) GetTrades(currencyPair exchange.CurrencyPair, since int64) ([]exchange.Trade, error) {
	panic("not implement")
}

//交易符号
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

func (fc *FCoin) GetTradeSymbols(currencyPair exchange.CurrencyPair) (*TradeSymbol, error) {
	if len(fc.tradeSymbols) == 0 {
		var err error
		fc.tradeSymbols, err = fc.getTradeSymbols()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range fc.tradeSymbols {
		if v.Name == strings.ToLower(currencyPair.ToSymbol("")) {
			return &fc.tradeSymbols[k], nil
		}
	}
	return nil, errors.New("symbol not found")
}
