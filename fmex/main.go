package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type TokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"data"`
}

type BuyRequest struct {
	ID              string `json:"id"`
	Amount          string `json:"amount"`
	PaymentCurrency string `json:"payment_currency"`
	Token           string `json:"token"`
}

var wg sync.WaitGroup

const (
	numGoroutines = 2
	numTokens     = 100

	OK    = "ok"
	ERROR = "error"

	contentType = "application/json;charset=utf-8"
	cookie      = "_ga=GA1.2.126722881.1551425649; __zlcmid=r6i3N3SCymOiYh; sandbox-token=0NL-A5jVNSWpHZG97U8gMom0jdY3iCryYR3b3mV6pZpx; _gid=GA1.2.549695496.1560576700; prd-token=_VgV0Uu5h2MrLzNy67AoHjlAkhkoLwSKcOwEDll1x082; LOGIN_FRESH=1; LOGIN_USERNAME=Jiangew"
	origin      = "https://www.fcoin.com"
	referer     = "https://www.fcoin.com/fmex"

	reqId     = "zj8rKpn7sAjhOYT9j8tgdQ"
	amount    = "12000"
	currency  = "usdt"
	prefixUrl = "https://www.fcoin.com/openapi/auth/v1/lightning_deals/"
	tokenUrl  = prefixUrl + reqId + "/token"
	buyUrl    = prefixUrl + reqId + "/buy"
)

func saveToken(tokens chan string) {
	defer wg.Done()
	client := &http.Client{}

	for {
		tokenReq, _ := http.NewRequest("POST", tokenUrl, nil)
		tokenReq.Header.Set("content-type", contentType)
		tokenReq.Header.Set("cookie", cookie)
		tokenReq.Header.Set("origin", origin)
		tokenReq.Header.Set("referer", referer)

		tokenResp, _ := client.Do(tokenReq)
		defer tokenResp.Body.Close()
		data, err := ioutil.ReadAll(tokenResp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}

		var tokenObject TokenResponse
		_ = json.Unmarshal(data, &tokenObject)
		if tokenObject.Status != OK {
			continue
		}
		tokens <- tokenObject.Token
		fmt.Println(tokenObject.Token)
	}
}

func buy(tokens chan string) {
	defer wg.Done()
	client := &http.Client{}

	for {
		token, _ := <-tokens

		buyReqBody := BuyRequest{
			ID:              reqId,
			Amount:          amount,
			PaymentCurrency: currency,
			Token:           token,
		}
		buyReqBytes, _ := json.Marshal(buyReqBody)

		buyReq, _ := http.NewRequest("POST", buyUrl, bytes.NewBuffer(buyReqBytes))
		buyReq.Header.Set("content-type", contentType)
		buyReq.Header.Set("cookie", cookie)
		buyReq.Header.Set("origin", origin)
		buyReq.Header.Set("referer", referer)

		buyResp, _ := client.Do(buyReq)
		defer buyResp.Body.Close()
		buyRespData, err := ioutil.ReadAll(buyResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buyRespData))
	}
}

func main() {
	tokens := make(chan string, numTokens)

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(numGoroutines)

	for count := 0; count < numGoroutines/2; count++ {
		go saveToken(tokens)
		go buy(tokens)
	}

	wg.Wait()
}
