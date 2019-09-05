package exchange

import "fmt"

func FmtBalance(balance float64, usdt float64, usdtFrozen float64, currency float64, currencyFrozen float64) string {
	return fmt.Sprintf("balance: %s, usdt: %s, usdtFrozen: %s, currency: %s, currencyFrozen: %s",
		FloatToStringForEx(balance),
		FloatToStringForEx(usdt),
		FloatToStringForEx(usdtFrozen),
		FloatToStringForEx(currency),
		FloatToStringForEx(currencyFrozen),
	)
}

func FmtBalanceExt(balance float64, usdt float64, usdtFrozen float64, fmex float64, fmexFrozen float64, ft float64, ftFrozen float64) string {
	return fmt.Sprintf("balance: %s, usdt: %s, usdtFrozen: %s, fmex: %s, fmexFrozen: %s, ft: %s, ftFrozen: %s",
		FloatToStringForEx(balance),
		FloatToStringForEx(usdt),
		FloatToStringForEx(usdtFrozen),
		FloatToStringForEx(fmex),
		FloatToStringForEx(fmexFrozen),
		FloatToStringForEx(ft),
		FloatToStringForEx(ftFrozen),
	)
}

func FmtOrder(symbol string, price float64, amount float64, state string, filledAmount float64) string {
	return fmt.Sprintf("{ symbol: %s, price: %s, amount: %s, state: %s, filledAmount: %s }",
		symbol,
		FloatToStringForEx(price),
		FloatToStringForEx(amount),
		state,
		FloatToStringForEx(filledAmount),
	)
}

func FmtTicker(ticker *Ticker) string {
	return fmt.Sprintf("symbol: %s, last: %s, lastVol: %s, buy: %s, buyVol: %s, sell: %s, sellVol: %s, high: %s, low: %s, baseVol: %s",
		ticker.Symbol,
		FloatToStringForEx(ticker.Last),
		FloatToStringForEx(ticker.LastVol),
		FloatToStringForEx(ticker.Buy),
		FloatToStringForEx(ticker.BuyVol),
		FloatToStringForEx(ticker.Sell),
		FloatToStringForEx(ticker.SellVol),
		FloatToStringForEx(ticker.High),
		FloatToStringForEx(ticker.Low),
		FloatToStringForEx(ticker.BaseVol),
	)
}

func FmtPaxMemoryStates(curBuyPrice float64, curSellPrice float64, lastBuyPrice float64, lastSellPrice float64, maxBuyPrice float64, minSellPrice float64) string {
	return fmt.Sprintf("curBuyPrice: %s, curSellPrice: %s, lastBuyPrice: %s, lastSellPrice: %s, maxBuyPrice: %s, minSellPrice: %s",
		FloatToStringForEx(curBuyPrice),
		FloatToStringForEx(curSellPrice),
		FloatToStringForEx(lastBuyPrice),
		FloatToStringForEx(lastSellPrice),
		FloatToStringForEx(maxBuyPrice),
		FloatToStringForEx(minSellPrice),
	)
}
