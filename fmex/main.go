package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jiangew/belex/exchange"
	"github.com/jiangew/belex/fcoin"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"time"
)

const (
	proxy  = "socks5://127.0.0.1:1086"
	key    = "30c2a9a181b945afaa7038f5326a2311"
	secret = "60078bb694424fc481590b4231b3e248"
	bot    = "941710186:AAFS--hy65duDvrBhPQisWIVgVCj8FMixcc"
)

var (
	symbol        = exchange.FMEX_USDT
	baseCurrency  = exchange.USDT
	quoteCurrency = exchange.FMEX
)

func main() {
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy(proxy)
	apiBuilder := fcoin.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey(key).APISecretkey(secret).Build(exchange.FCOIN)

	bot, err := tgbotapi.NewBotAPI(bot)
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// async send telegram message
	go sendMessage(api, bot, updates)
}

func sendMessage(api exchange.API, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "b":
			usdtAccount, _ := api.GetSubAccount(baseCurrency)
			currencyAccount, _ := api.GetSubAccount(quoteCurrency)
			taker, _ := api.GetTicker(symbol)
			currencyToUsdt := decimal.NewFromFloat(currencyAccount.Balance).Mul(decimal.NewFromFloat(taker.Sell))
			balance := decimal.NewFromFloat(usdtAccount.Balance).Add(currencyToUsdt)
			balanceOut, _ := strconv.ParseFloat(balance.String(), 64)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("balance: %s, usdt: %s, usdtFrozen: %s, currency: %s, currencyFrozen: %s",
				fmt.Sprintf("%.4f", balanceOut),
				fmt.Sprintf("%.4f", usdtAccount.Available),
				fmt.Sprintf("%.4f", usdtAccount.Frozen),
				fmt.Sprintf("%.4f", currencyAccount.Available),
				fmt.Sprintf("%.4f", currencyAccount.Frozen),
			))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "o":
			orders, _ := api.GetActiveOrders(symbol)
			buyCount := 0
			var buyOrders []exchange.NewOrder
			sellCount := 0
			var sellOrders []exchange.NewOrder
			if len(orders) > 0 {
				for _, order := range orders {
					if order.Side == "buy" {
						buyCount++
						buyOrders = append(buyOrders, order)
					} else if order.Side == "sell" {
						sellCount++
						sellOrders = append(sellOrders, order)
					}
				}
			}
			msgBody := ""
			if len(orders) > 0 {
				buyBytes, _ := json.Marshal(buyOrders)
				sellBytes, _ := json.Marshal(sellOrders)
				msgBody = fmt.Sprintf("buyCount: %d, buyOrders: %s, sellCount: %d, sellOrders: %s", buyCount, string(buyBytes), sellCount, string(sellBytes))
			} else {
				msgBody = "there is no active orders."
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "t":
			taker, _ := api.GetTicker(symbol)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("symbol: %s, last: %s, lastVol: %s, buy: %s, buyVol: %s, sell: %s, sellVol: %s, high: %s, low: %s, baseVol: %s",
				taker.Symbol,
				fmt.Sprintf("%.4f", taker.Last),
				fmt.Sprintf("%.4f", taker.LastVol),
				fmt.Sprintf("%.4f", taker.Buy),
				fmt.Sprintf("%.4f", taker.BuyVol),
				fmt.Sprintf("%.4f", taker.Sell),
				fmt.Sprintf("%.4f", taker.SellVol),
				fmt.Sprintf("%.4f", taker.High),
				fmt.Sprintf("%.4f", taker.Low),
				fmt.Sprintf("%.4f", taker.BaseVol),
			))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
