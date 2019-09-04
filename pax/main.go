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
	key    = "d314f40f344646148246e1b314df721c"
	secret = "1869f3b5da2b4b0bb99cf8a9b692e828"
	bot    = "960133387:AAGZ3dZ1FPO-lVJmVTUYsMxDZFUR5WDEEc0"
)

var (
	symbol        = exchange.PAX_USDT
	baseCurrency  = exchange.USDT
	quoteCurrency = exchange.PAX

	availableLimit = float64(300)

	upRise   = float64(10003.0 / 10000.0)
	downRise = float64(9997.0 / 10000.0)

	curBuyPrice   = float64(0)
	curSellPrice  = float64(0)
	lastBuyPrice  = float64(0)
	lastSellPrice = float64(0)
	maxBuyPrice   = float64(0)
	minSellPrice  = float64(0)
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

	// cancel orders when start
	orders, _ := api.GetActiveOrders(symbol)
	if len(orders) > 0 {
		for _, order := range orders {
			cancel, _ := api.CancelOrder(order.ID, symbol)
			log.Println("cancel order when start:", order.ID, "ret:", cancel)
		}
	}

	// main quant
	for {
		ticker, err := api.GetTicker(symbol)
		if err != nil {
			log.Println("usdt account got error:", err)
			continue
		} else {
			curBuyPrice = ticker.Buy
			curSellPrice = ticker.Sell
			log.Println("bid price:", curBuyPrice)
			log.Println("ask price:", curSellPrice)
		}

		orders, _ := api.GetActiveOrders(symbol)
		if len(orders) > 0 {
			for _, order := range orders {
				if order.Side == "buy" {
					if order.Price != curBuyPrice {
						cancel, _ := api.CancelOrder(order.ID, symbol)
						log.Println("cancel buy:", order.ID, "ret:", cancel)
					}
				} else if order.Side == "sell" {
					if order.Price != curSellPrice {
						cancel, _ := api.CancelOrder(order.ID, symbol)
						log.Println("cancel sell:", order.ID, "ret:", cancel)
					}
				}
			}
		}

		usdtAccount, err := api.GetSubAccount(baseCurrency)
		if err != nil {
			log.Println("usdt account got error:", err)
		} else {
			if usdtAccount.Available > availableLimit {
				if (maxBuyPrice > 0 && curBuyPrice > maxBuyPrice) {
					log.Println("limit buy exceeded limit price:", curBuyPrice)
				} else {
					isOrderable, _ := api.IsOrderable(symbol)
					if isOrderable {
						amount := (usdtAccount.Available - 1) / curBuyPrice
						if amount > 1 {
							buyOrder, err := api.LimitBuy(exchange.FloatToStringForEx(amount), exchange.FloatToStringForEx(curBuyPrice), symbol)
							if err != nil {
								log.Println("limit buy amount:", amount, "price:", curBuyPrice, "error:", err)
							} else {
								log.Println("limit buy amount:", amount, "price:", curBuyPrice, "success:", buyOrder.ID)
								lastBuyPrice = curBuyPrice
								maxBuyPrice = curBuyPrice * upRise
								minSellPrice = curBuyPrice * downRise
							}
						}
					} else {
						log.Println("limit buy isOrderable:", false)
					}
				}
			}
		}

		currencyAccount, err := api.GetSubAccount(quoteCurrency)
		if err != nil {
			log.Println("currency account got error:", err)
		} else {
			if currencyAccount.Available > availableLimit {
				if (minSellPrice > 0 && curSellPrice < minSellPrice) {
					log.Println("limit sell exceeded limit price:", curSellPrice)
				} else {
					isOrderable, _ := api.IsOrderable(symbol)
					if isOrderable {
						amount := currencyAccount.Available - 1
						if amount > 1 {
							sellOrder, err := api.LimitSell(exchange.FloatToStringForEx(amount), exchange.FloatToStringForEx(curSellPrice), symbol)
							if err != nil {
								log.Println("limit sell amount:", amount, "price:", curSellPrice, "error:", err)
							} else {
								log.Println("limit sell amount:", amount, "price:", curSellPrice, "success:", sellOrder.ID)
								lastSellPrice = curSellPrice
								minSellPrice = curSellPrice * downRise
								maxBuyPrice = curSellPrice * upRise
							}
						}
					} else {
						log.Println("limit sell isOrderable:", false)
					}
				}
			}
		}

		time.Sleep(250 * time.Millisecond)
	}
}

func sendMessage(api exchange.API, bot *tgbot.BotAPI, updates tgbot.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "func":
			msgBody := fmt.Sprintf("b -> %s\n, o -> %s\n, t -> %s\n, m -> %s\n, start -> %s\n, stop -> %s\n, fb -> %s\n, fo -> %s\n, fbo -> %s\n, fso -> %s\n",
				"pax balance",
				"pax stats orders",
				"pax ticker",
				"pax exchange states in memory",
				"start service",
				"stop service",
				"ft balance",
				"ft stats orders",
				"ft buy orders",
				"ft sell orders",
			)
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
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
		case "t":
			ticker, _ := api.GetTicker(symbol)
			msg := tgbot.NewMessage(update.Message.Chat.ID, exchange.FmtTicker(ticker))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "m":
			msgBody := exchange.FmtPaxMemoryStates(curBuyPrice, curSellPrice, lastBuyPrice, lastSellPrice, maxBuyPrice, minSellPrice)
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "start":
			maxBuyPrice = float64(0)
			minSellPrice = float64(0)
			msg := tgbot.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been cleared.")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "stop":
			ticker, _ := api.GetTicker(symbol)
			maxBuyPrice = ticker.Buy / 2
			minSellPrice = ticker.Sell * 2

			orders, _ := api.GetActiveOrders(symbol)
			if len(orders) > 0 {
				for _, order := range orders {
					cancel, _ := api.CancelOrder(order.ID, symbol)
					log.Println("cancel order when stop:", order.ID, "ret:", cancel)
				}
			}

			msg := tgbot.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been set.")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "fb":
			ftAccount, _ := api.GetSubAccount(exchange.FT)
			ticker, _ := api.GetTicker(exchange.FT_USDT)
			ftToUsdt := decimal.NewFromFloat(ftAccount.Balance).Mul(decimal.NewFromFloat(ticker.Sell))
			balanceOut, _ := strconv.ParseFloat(ftToUsdt.String(), 64)

			msgBody := exchange.FmtCurrencyBalance(balanceOut, ftAccount.Available, ftAccount.Frozen)
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "fo":
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
		case "fbo":
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
				msgBody = fmt.Sprintf("buyCount: %d\n, buyOrders: %s", buyCount, strings.Join(buyOrders, ",\n"))
			} else {
				msgBody = "there is no buy active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "fso":
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
				msgBody = fmt.Sprintf("sellCount: %d\n, sellOrders: %s", sellCount, strings.Join(sellOrders, ",\n"))
			} else {
				msgBody = "there is no sell active orders."
			}
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
