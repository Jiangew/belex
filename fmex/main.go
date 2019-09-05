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
		case "x":
			msgBody := "b -> account balance\n" +
				"\n" +
				"mt -> fmex ticker\n" +
				"mo -> fmex stats orders\n" +
				"mb -> fmex buy orders\n" +
				"ms -> fmex sell orders\n" +
				"mc -> fmex cancel orders\n" +
				"\n" +
				"tt -> ft ticker\n" +
				"to -> ft stats orders\n" +
				"tb -> ft buy orders\n" +
				"ts -> ft sell orders\n" +
				"tc -> ft cancel orders";
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "b":
			usdtAccount, _ := api.GetSubAccount(exchange.USDT)

			fmexAccount, _ := api.GetSubAccount(exchange.FMEX)
			fmexTicker, _ := api.GetTicker(exchange.FMEX_USDT)
			fmexToUsdt := decimal.NewFromFloat(fmexAccount.Balance).Mul(decimal.NewFromFloat(fmexTicker.Sell))

			ftAccount, _ := api.GetSubAccount(exchange.FT)
			ftTicker, _ := api.GetTicker(exchange.FT_USDT)
			ftToUsdt := decimal.NewFromFloat(ftAccount.Balance).Mul(decimal.NewFromFloat(ftTicker.Sell))

			balance := decimal.NewFromFloat(usdtAccount.Balance).Add(fmexToUsdt).Add(ftToUsdt)
			balanceOut, _ := strconv.ParseFloat(balance.String(), 64)

			msgBody := exchange.FmtBalanceExt(balanceOut, usdtAccount.Available, usdtAccount.Frozen, fmexAccount.Available, fmexAccount.Frozen, ftAccount.Available, ftAccount.Frozen)
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "mt":
			ticker, _ := api.GetTicker(exchange.FMEX_USDT)
			msg := tgbot.NewMessage(update.Message.Chat.ID, exchange.FmtTicker(ticker))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "mo":
			orders, _ := api.GetActiveOrders(exchange.FMEX_USDT)
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
		case "mb":
			orders, _ := api.GetActiveOrders(exchange.FMEX_USDT)
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
				msgBody = fmt.Sprintf("buyCount: %d\nbuyOrders: %s", buyCount, strings.Join(buyOrders, ",\n"))
			} else {
				msgBody = "there is no buy active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "ms":
			orders, _ := api.GetActiveOrders(exchange.FMEX_USDT)
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
				msgBody = fmt.Sprintf("sellCount: %d\nsellOrders: %s", sellCount, strings.Join(sellOrders, ",\n"))
			} else {
				msgBody = "there is no sell active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "mc":
			orders, _ := api.GetActiveOrders(exchange.FMEX_USDT)
			if len(orders) > 0 {
				for _, order := range orders {
					cancel, _ := api.CancelOrder(order.ID, exchange.FMEX_USDT)
					log.Println("cancel order:", order.ID, "ret:", cancel)
				}
			}

			msgBody := fmt.Sprintf("symbol: %s count: %d active orders has been canceled.", exchange.FMEX_USDT.String(), len(orders))
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "tt":
			ticker, _ := api.GetTicker(exchange.FT_USDT)
			msg := tgbot.NewMessage(update.Message.Chat.ID, exchange.FmtTicker(ticker))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "to":
			orders, _ := api.GetActiveOrders(exchange.FT_USDT)
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
		case "tb":
			orders, _ := api.GetActiveOrders(exchange.FT_USDT)
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
				msgBody = fmt.Sprintf("buyCount: %d\nbuyOrders: %s", buyCount, strings.Join(buyOrders, ",\n"))
			} else {
				msgBody = "there is no buy active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "ts":
			orders, _ := api.GetActiveOrders(exchange.FT_USDT)
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
				msgBody = fmt.Sprintf("sellCount: %d\nsellOrders: %s", sellCount, strings.Join(sellOrders, ",\n"))
			} else {
				msgBody = "there is no sell active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "tc":
			orders, _ := api.GetActiveOrders(exchange.FT_USDT)
			if len(orders) > 0 {
				for _, order := range orders {
					cancel, _ := api.CancelOrder(order.ID, exchange.FT_USDT)
					log.Println("cancel order:", order.ID, "ret:", cancel)
				}
			}

			msgBody := fmt.Sprintf("symbol: %s count: %d active orders has been canceled.", exchange.FT_USDT.String(), len(orders))
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
