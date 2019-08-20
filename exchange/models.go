package exchange

import (
	"time"
)

type NewOrder struct {
	ID           string
	Symbol       string
	Side         string
	OrderType    string
	Price        float64
	Amount       float64
	State        string
	FilledAmount float64
	FillFee      float64
	Source       string
	CreatedAt    uint64
}

type Order struct {
	Price      float64
	Amount     float64
	AvgPrice   float64
	DealAmount float64
	Fee        float64
	Cid        string //clientId
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
	Last    float64 `json:"last,string"`     //最新成交价
	LastVol float64 `json:"last_vol,string"` //最近一笔成交的成交量
	Buy     float64 `json:"buy,string"`      //最大买一价
	BuyVol  float64 `json:"buy_vol,string"`  //最大买一量
	Sell    float64 `json:"sell,string"`     //最小卖一价
	SellVol float64 `json:"sell_vol,string"` //最小卖一量
	High    float64 `json:"high,string"`     //24小时内最高价
	Low     float64 `json:"low,string"`      //24小时内最低价
	BaseVol float64 `json:"base_vol,string"` //24小时内基准货币成交量
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
	Symbol  string
	UTime   time.Time
	AskList DepthRecords //Descending order
	BidList DepthRecords //Descending order
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
