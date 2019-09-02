package main

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jiangew/belex/exchange"
	"github.com/jiangew/belex/fcoin"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"strings"
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

	bot, err := tgbot.NewBotAPI(bot)
	if err != nil {
		log.Panic(err)
	}
	u := tgbot.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// async send telegram message
	go sendMessage(api, bot, updates)

	for {
		time.Sleep(250 * time.Millisecond)
	}
}

func sendMessage(api exchange.API, bot *tgbot.BotAPI, updates tgbot.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "b":
			usdtAccount, _ := api.GetSubAccount(baseCurrency)
			currencyAccount, _ := api.GetSubAccount(quoteCurrency)
			ticker, _ := api.GetTicker(symbol)
			currencyToUsdt := decimal.NewFromFloat(currencyAccount.Balance).Mul(decimal.NewFromFloat(ticker.Sell))
			balance := decimal.NewFromFloat(usdtAccount.Balance).Add(currencyToUsdt)
			balanceOut, _ := strconv.ParseFloat(balance.String(), 64)

			msgBody := exchange.FmtBalance(balanceOut, usdtAccount.Available, usdtAccount.Frozen, currencyAccount.Available, currencyAccount.Frozen)
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "o":
			orders, _ := api.GetActiveOrders(symbol)
			buyCount := 0
			sellCount := 0
			if len(orders) > 0 {
				for _, order := range orders {
					if order.Side == "buy" {
						buyCount++
					} else if order.Side == "sell" {
						sellCount++
					}
				}
			}

			msgBody := ""
			if len(orders) > 0 {
				msgBody = fmt.Sprintf("buyCount: %d, sellCount: %d", buyCount, sellCount)
			} else {
				msgBody = "there is no active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "ob":
			orders, _ := api.GetActiveOrders(symbol)
			buyCount := 0
			var buyOrders []string
			if len(orders) > 0 {
				for _, order := range orders {
					ord := exchange.FmtOrder(order.Symbol, order.Price, order.Amount, order.State, order.FilledAmount)
					if order.Side == "buy" {
						buyCount++
						buyOrders = append(buyOrders, ord)
					}
				}
			}

			msgBody := ""
			if len(orders) > 0 {
				msgBody = fmt.Sprintf("buyCount: %d, buyOrders: %s", buyCount, strings.Join(buyOrders, ", "))
			} else {
				msgBody = "there is no buy active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "os":
			orders, _ := api.GetActiveOrders(symbol)
			sellCount := 0
			var sellOrders []string
			if len(orders) > 0 {
				for _, order := range orders {
					ord := exchange.FmtOrder(order.Symbol, order.Price, order.Amount, order.State, order.FilledAmount)
					if order.Side == "sell" {
						sellCount++
						sellOrders = append(sellOrders, ord)
					}
				}
			}

			msgBody := ""
			if len(orders) > 0 {
				msgBody = fmt.Sprintf("sellCount: %d, sellOrders: %s", sellCount, strings.Join(sellOrders, ", "))
			} else {
				msgBody = "there is no sell active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "t":
			ticker, _ := api.GetTicker(symbol)
			msg := tgbot.NewMessage(update.Message.Chat.ID, exchange.FmtTicker(ticker))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "c":
			orders, _ := api.GetActiveOrders(symbol)
			if len(orders) > 0 {
				for _, order := range orders {
					cancel, _ := api.CancelOrder(order.ID, symbol)
					log.Println("cancel order:", order.ID, "ret:", cancel)
				}
			}

			msgBody := fmt.Sprintf("symbol: %s count: %d active orders has been canceled.", symbol.String(), len(orders))
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
