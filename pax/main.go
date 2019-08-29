package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jiangew/belex/exchange"
	"github.com/jiangew/belex/fcoin"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"time"
)

//var wg sync.WaitGroup

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

	bot, err := tgbotapi.NewBotAPI(bot)
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
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
		taker, err := api.GetTicker(symbol)
		if err != nil {
			log.Println("usdt account got error:", err)
			continue
		} else {
			curBuyPrice = taker.Buy
			curSellPrice = taker.Sell
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
							buyOrder, err := api.LimitBuy(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", curBuyPrice), symbol)
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
							sellOrder, err := api.LimitSell(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", curSellPrice), symbol)
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgBody)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "m":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("curBuyPrice: %s, curSellPrice: %s, lastBuyPrice: %s, lastSellPrice: %s, maxBuyPrice: %s, minSellPrice: %s",
				fmt.Sprintf("%.4f", curBuyPrice),
				fmt.Sprintf("%.4f", curSellPrice),
				fmt.Sprintf("%.4f", lastBuyPrice),
				fmt.Sprintf("%.4f", lastSellPrice),
				fmt.Sprintf("%.4f", maxBuyPrice),
				fmt.Sprintf("%.4f", minSellPrice),
			))
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
		case "start":
			maxBuyPrice = float64(0)
			minSellPrice = float64(0)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been cleared.")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case "stop":
			taker, _ := api.GetTicker(symbol)
			maxBuyPrice = taker.Buy / 2
			minSellPrice = taker.Sell * 2

			orders, _ := api.GetActiveOrders(symbol)
			if len(orders) > 0 {
				for _, order := range orders {
					cancel, _ := api.CancelOrder(order.ID, symbol)
					log.Println("cancel order when stop:", order.ID, "ret:", cancel)
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "max buy and min sell limit price in memory has been set.")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		}
	}
}
