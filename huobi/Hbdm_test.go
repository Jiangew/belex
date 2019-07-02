package huobi

import (
	"github.com/jiangew/belex"
	"testing"
	"time"
)

var dm = NewHbdm(&belex.APIConfig{
	Endpoint:     "https://api.hbdm.com",
	HttpClient:   httpProxyClient,
	ApiKey:       "",
	ApiSecretKey: ""})

func TestHbdm_GetFutureUserinfo(t *testing.T) {
	t.Log(dm.GetFutureUserinfo())
}

func TestHbdm_GetFuturePosition(t *testing.T) {
	t.Log(dm.GetFuturePosition(belex.BTC_USD, belex.QUARTER_CONTRACT))
}

func TestHbdm_PlaceFutureOrder(t *testing.T) {
	t.Log(dm.PlaceFutureOrder(belex.BTC_USD, belex.QUARTER_CONTRACT, "3800", "1", belex.OPEN_BUY, 0, 20))
}

func TestHbdm_FutureCancelOrder(t *testing.T) {
	t.Log(dm.FutureCancelOrder(belex.BTC_USD, belex.QUARTER_CONTRACT, "6"))
}

func TestHbdm_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(dm.GetUnfinishFutureOrders(belex.BTC_USD, belex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureOrders(t *testing.T) {
	t.Log(dm.GetFutureOrders([]string{"6", "5"}, belex.BTC_USD, belex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureOrder(t *testing.T) {
	t.Log(dm.GetFutureOrder("6", belex.BTC_USD, belex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureTicker(t *testing.T) {
	t.Log(dm.GetFutureTicker(belex.EOS_USD, belex.QUARTER_CONTRACT))
}

func TestHbdm_GetFutureDepth(t *testing.T) {
	dep, err := dm.GetFutureDepth(belex.BTC_USD, belex.QUARTER_CONTRACT, 0)
	t.Log(err)
	t.Logf("%+v\n%+v", dep.AskList, dep.BidList)
}
func TestHbdm_GetFutureIndex(t *testing.T) {
	t.Log(dm.GetFutureIndex(belex.BTC_USD))
}

func TestHbdm_GetFutureEstimatedPrice(t *testing.T) {
	t.Log(dm.GetFutureEstimatedPrice(belex.BTC_USD))
}

func TestHbdm_GetKlineRecords(t *testing.T) {
	klines, _ := dm.GetKlineRecords(belex.QUARTER_CONTRACT, belex.EOS_USD, belex.KLINE_PERIOD_1MIN, 20, 0)
	for _, k := range klines {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(k.Pair, tt, k.Open, k.Close, k.High, k.Low, k.Vol, k.Vol2)
	}
}
