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

//var wg sync.WaitGroup

func main() {
	//apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1086")
	apiBuilder := fcoin.NewAPIBuilder().HttpTimeout(5 * time.Second)
	api := apiBuilder.APIKey("1412ac27e3f741c796f7c4600069d9f1").APISecretkey("4843754749be46919d986142917f06d7").Build(exchange.FCOIN)

	bot, err := tgbotapi.NewBotAPI("960133387:AAGZ3dZ1FPO-lVJmVTUYsMxDZFUR5WDEEc0")
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Text {
			case "b":
				usdtAccount, _ := api.GetSubAccount(exchange.USDT)
				currencyAccount, _ := api.GetSubAccount(exchange.PAX)
				taker, _ := api.GetTicker(exchange.PAX_USDT)

				currencyToUsdt := decimal.NewFromFloat(currencyAccount.Balance).Mul(decimal.NewFromFloat(taker.Sell))
				balance := decimal.NewFromFloat(usdtAccount.Balance).Add(currencyToUsdt)
				balanceOut, _ := strconv.ParseFloat(balance.String(), 64)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("usdt: %s, usdtFrozen: %s, currency: %s, currencyFrozen: %s, balance: %s",
					fmt.Sprintf("%.4f", usdtAccount.Available),
					fmt.Sprintf("%.4f", usdtAccount.Frozen),
					fmt.Sprintf("%.4f", currencyAccount.Available),
					fmt.Sprintf("%.4f", currencyAccount.Frozen),
					fmt.Sprintf("%.4f", balanceOut)))
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			case "t":
				taker, _ := api.GetTicker(exchange.PAX_USDT)
				takerBytes, _ := json.Marshal(taker)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(takerBytes))
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			case "o":
				orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
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
					msgBody = fmt.Sprintf("orderCount: %d, buyCount: %d, sellCount: %d", len(orders), buyCount, sellCount)
				} else {
					msgBody = "there is no active orders."
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgBody)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}
		}
	}()

	orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
	if len(orders) > 0 {
		for _, order := range orders {
			cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
			log.Println("cancel order:", order.ID, "ret:", cancel)
		}
	}

	buyPrice := float64(0)
	maxBuyPrice := float64(0)
	sellPrice := float64(0)
	minSellPrice := float64(0)

	for {
		taker, err := api.GetTicker(exchange.PAX_USDT)
		if err != nil {
			log.Println("usdt account got error:", err)
			continue
		} else {
			buyPrice = taker.Buy
			sellPrice = taker.Sell
			log.Println("bid price:", buyPrice)
			log.Println("ask price:", sellPrice)
		}

		orders, _ := api.GetActiveOrders(exchange.PAX_USDT)
		if len(orders) > 0 {
			for _, order := range orders {
				if order.Side == "buy" {
					if order.Price != buyPrice {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("cancel buy:", order.ID, "ret:", cancel)
					}
				} else if order.Side == "sell" {
					if order.Price != sellPrice {
						cancel, _ := api.CancelOrder(order.ID, exchange.PAX_USDT)
						log.Println("cancel sell:", order.ID, "ret:", cancel)
					}
				}
			}
		}

		usdtAccount, err := api.GetSubAccount(exchange.USDT)
		if err != nil {
			log.Println("usdt account got error:", err)
		} else {
			if usdtAccount.Available > 300 {
				if (maxBuyPrice > 0 && buyPrice > maxBuyPrice) {
					log.Println("limit buy exceeded limit price:", buyPrice)
				} else {
					isOrderable, _ := api.IsOrderable(exchange.PAX_USDT)
					if isOrderable {
						amount := (usdtAccount.Available - 1) / buyPrice
						if amount > 1 {
							buyOrder, err := api.LimitBuy(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", buyPrice), exchange.PAX_USDT)
							if err != nil {
								log.Println("limit buy amount:", amount, "price:", buyPrice, "error:", err)
							} else {
								log.Println("limit buy amount:", amount, "price:", buyPrice, "success:", buyOrder.ID)
								maxBuyPrice = buyPrice * 10005 / 10000
							}
						}
					} else {
						log.Println("limit buy isOrderable:", false)
					}
				}
			}
		}

		currencyAccount, err := api.GetSubAccount(exchange.PAX)
		if err != nil {
			log.Println("pax account got error:", err)
		} else {
			if currencyAccount.Available > 300 {
				if (minSellPrice > 0 && sellPrice < minSellPrice) {
					log.Println("limit sell exceeded limit price:", sellPrice)
				} else {
					isOrderable, _ := api.IsOrderable(exchange.PAX_USDT)
					if isOrderable {
						amount := currencyAccount.Available - 1
						if amount > 1 {
							sellOrder, err := api.LimitSell(fmt.Sprintf("%.4f", amount), fmt.Sprintf("%.4f", sellPrice), exchange.PAX_USDT)
							if err != nil {
								log.Println("limit sell amount:", amount, "price:", sellPrice, "error:", err)
							} else {
								log.Println("limit sell amount:", amount, "price:", sellPrice, "success:", sellOrder.ID)
								minSellPrice = sellPrice * 9995 / 10000
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
