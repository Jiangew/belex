package main

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jiangew/belex/exchange"
	"github.com/jiangew/belex/fcoin"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
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
				isOrderable, _ := api.IsOrderable(symbol)
				if isOrderable {
					if (maxBuyPrice > 0 && curBuyPrice > maxBuyPrice) {
						log.Println("limit buy exceeded limit price:", curBuyPrice)
					} else {
						amount := (usdtAccount.Available - 1) / curBuyPrice
						if amount > 1 {
							buyOrder, err := api.LimitBuy(exchange.FloatToStringForEx(amount), exchange.FloatToStringForEx(curBuyPrice), symbol)
							if err != nil {
								log.Println("limit buy amount:", amount, "price:", curBuyPrice, "error:", err)
							} else {
								log.Println("limit buy amount:", amount, "price:", curBuyPrice, "success:", buyOrder.ID)
								lastBuyPrice = curBuyPrice
								minSellPrice = curBuyPrice * downRise
							}
						}
					}
				} else {
					log.Println("limit buy isOrderable:", false)
				}
			}
		}

		currencyAccount, err := api.GetSubAccount(quoteCurrency)
		if err != nil {
			log.Println("currency account got error:", err)
		} else {
			if currencyAccount.Available > availableLimit {
				isOrderable, _ := api.IsOrderable(symbol)
				if isOrderable {
					if (minSellPrice > 0 && curSellPrice < minSellPrice) {
						log.Println("limit sell exceeded limit price:", curSellPrice)
					} else {
						amount := currencyAccount.Available - 1
						if amount > 1 {
							sellOrder, err := api.LimitSell(exchange.FloatToStringForEx(amount), exchange.FloatToStringForEx(curSellPrice), symbol)
							if err != nil {
								log.Println("limit sell amount:", amount, "price:", curSellPrice, "error:", err)
							} else {
								log.Println("limit sell amount:", amount, "price:", curSellPrice, "success:", sellOrder.ID)
								lastSellPrice = curSellPrice
								maxBuyPrice = curSellPrice * upRise
							}
						}
					}
				} else {
					log.Println("limit sell isOrderable:", false)
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
		case "x":
			msgBody := "b -> stats balance\n" +
				"o -> stats orders\n" +
				"t -> ticker\n" +
				"m -> exchange states in memory\n" +
				"start -> start service\n" +
				"stop -> stop service";
			msg := tgbot.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "b":
			usdtAccount, _ := api.GetSubAccount(baseCurrency)
			ftAccount, _ := api.GetSubAccount(exchange.FT)

			currencyAccount, _ := api.GetSubAccount(quoteCurrency)
			currencyTicker, _ := api.GetTicker(symbol)
			currencyToUsdt := decimal.NewFromFloat(currencyAccount.Balance).Mul(decimal.NewFromFloat(currencyTicker.Sell))

			balance := decimal.NewFromFloat(usdtAccount.Balance).Add(currencyToUsdt)
			balanceOut, _ := strconv.ParseFloat(balance.String(), 64)

			msgBody := exchange.FmtBalance(balanceOut, usdtAccount.Available, usdtAccount.Frozen, currencyAccount.Available, currencyAccount.Frozen, ftAccount.Available, ftAccount.Frozen)
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
			msg := tgbot.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been set.")
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

			msg := tgbot.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been cleared.")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
