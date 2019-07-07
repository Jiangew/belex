package exchange

import (
	"net/http"
	"time"
)

type NewOrder struct {
	ID            string
	Symbol        string
	Side          string
	OrderType     string
	Price         float64
	Amount        float64
	State         string
	FilledAmount  float64
	FillFee       float64
	Source        string
	CreatedAt     uint64
}

type Order struct {
	Price      float64
	Amount     float64
	AvgPrice   float64
	DealAmount float64
	Fee        float64
	Cid        string //客户端自定义ID
	OrderID2   string
	OrderID    int //deprecated
	OrderTime  int
	Status     OrderState
	Symbol     Symbol
	Side       TradeSide
}

type Trade struct {
	Tid    int64     `json:"tid"`
	Type   TradeSide `json:"type"`
	Amount float64   `json:"amount,string"`
	Price  float64   `json:"price,string"`
	Date   int64     `json:"date_ms"`
	Symbol Symbol    `json:"omitempty"`
}

type SubAccount struct {
	Currency      string
	Available     float64
	Frozen        float64
	DemandDeposit float64
	LockDeposit   float64
	Balance       float64
}

type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[Currency]SubAccount
}

type Ticker struct {
	Symbol  string  `json:"symbol"`
	Date    uint64  `json:"date"`
	Last    float64 `json:"last,string"`
	LastVol float64 `json:"last_vol,string"`
	Buy     float64 `json:"buy,string"`
	BuyVol  float64 `json:"buy_vol,string"`
	Sell    float64 `json:"sell,string"`
	SellVol float64 `json:"sell_vol,string"`
	High    float64 `json:"high,string"`
	Low     float64 `json:"low,string"`
	Vol     float64 `json:"vol,string"`
}

type DepthRecord struct {
	Price  float64
	Amount float64
}

type DepthRecords []DepthRecord

func (dr DepthRecords) Len() int {
	return len(dr)
}

func (dr DepthRecords) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthRecords) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	Symbol       string
	UTime        time.Time
	AskList      DepthRecords //Descending order
	BidList      DepthRecords //Descending order
}

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net, huobi.pro

	Lever int //杠杆倍数, for future
}

type Kline struct {
	Symbol    Symbol
	Timestamp int64
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Vol       float64
}

type FutureKline struct {
	*Kline
	Vol2 float64 //个数
}

type FutureSubAccount struct {
	Currency      Currency
	AccountRights float64 //账户权益
	KeepDeposit   float64 //保证金
	ProfitReal    float64 //已实现盈亏
	ProfitUnreal  float64
	RiskRate      float64 //保证金率
}

type FutureAccount struct {
	FutureSubAccounts map[Currency]FutureSubAccount
}

type FutureOrder struct {
	OrderID2     string //请尽量用这个字段替代OrderID字段
	Price        float64
	Amount       float64
	AvgPrice     float64
	DealAmount   float64
	OrderID      int64 //deprecated
	OrderTime    int64
	Status       OrderState
	Currency     Symbol
	OType        int     //1：开多 2：开空 3：平多 4： 平空
	LeverRate    int     //倍数
	Fee          float64 //手续费
	ContractName string
}

type FuturePosition struct {
	BuyAmount      float64
	BuyAvailable   float64
	BuyPriceAvg    float64
	BuyPriceCost   float64
	BuyProfitReal  float64
	CreateDate     int64
	LeverRate      int
	SellAmount     float64
	SellAvailable  float64
	SellPriceAvg   float64
	SellPriceCost  float64
	SellProfitReal float64
	Symbol         Symbol //btc_usd:比特币,ltc_usd:莱特币
	ContractType   string
	ContractId     int64
	ForceLiquPrice float64 //预估爆仓价
}
