package api

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
	"strconv"
	"time"
)

// 交易网格
var gPriceToQuantityMap = map[string]string{
	"10":   "4",
	"10.2": "3.92",
	"10.4": "3.85",
	"10.6": "3.77",
	"10.8": "3.7",
	"11":   "3.64",
	"11.2": "3.57",
	"11.4": "3.51",
	"11.6": "3.45",
	"11.8": "3.39",
	"12":   "3.33",
	"12.2": "3.28",
	"12.4": "3.23",
	"12.6": "3.17",
	"12.8": "3.13",
	"13":   "3.08",
	"13.2": "3.03",
	"13.4": "2.99",
	"13.6": "2.94",
	"13.8": "2.9",
	"14":   "2.86",
	"14.2": "2.82",
	"14.4": "1.39",
	"14.6": "1.37",
	"14.8": "1.35",
	"15":   "1.33",
	"15.2": "1.32",
	"15.4": "1.3",
	"15.6": "1.28",
	"15.8": "1.27",
	"16":   "1.25",
	"16.2": "1.23",
	"16.4": "1.22",
	"16.6": "1.2",
	"16.8": "1.19",
	"17":   "1.18",
	"17.2": "1.16",
	"17.4": "1.15",
	"17.6": "1.14",
	"17.8": "1.12",
	"18":   "1.11",
}

// 价格列表
var gPriceList = []string{"10", "10.2", "10.4", "10.6", "10.8", "11", "11.2", "11.4", "11.6", "11.8", "12", "12.2", "12.4", "12.6", "12.8", "13", "13.2", "13.4", "13.6", "13.8", "14", "14.2", "14.4", "14.6", "14.8", "15", "15.2", "15.4", "15.6", "15.8", "16", "16.2", "16.4", "16.6", "16.8", "17", "17.2", "17.4", "17.6", "17.8", "18"}

func queryUnSettledTxs() (txs []*BnTxs, err error) {
	err = gDb.Raw("select * from bn_txs where status=0").Scan(&txs).Error
	if err != nil {
		fmt.Println("queryUnSettledTxs err", err)
	}
	return
}

func querySettledTxs() (txs []*BnTxs) {
	err := gDb.Raw("select * from bn_txs where status=1").Scan(&txs).Error
	if err != nil {
		fmt.Println("querySettledTxs err", err)
	}
	return
}

func getLeftAndRightPrice(price string) (sm, lg string) {
	for k := range gPriceList {
		if gPriceList[k] <= price && gPriceList[k+1] >= price {
			return gPriceList[k], gPriceList[k+1]
		}
	}
	panic("unexpected")
}

// 买的的时候挂委托单，卖的时候现价单，先这样做
func AutoGo(symbol string) {
	price := GetPrice(symbol)
	if price == "" {
		fmt.Println("AutoGo can not get latest price")
		return
	}
	smPri, _ := getLeftAndRightPrice(price)  // 14.8 => (14.6,14.8) 或者 14.81 => (14.8,15)
	_smPri, _ := getLeftAndRightPrice(smPri) // 14.6 => (14.4,16.6)
	txs, err := queryUnSettledTxs()
	if err != nil {
		return
	}
	fmt.Println("AutoGo unsettled order length: ", len(txs), "price:", price)
	checkBuy(txs, smPri, symbol)
	checkSell(txs, _smPri)
}

func CheckOrd() {
	// 处理卖出订单
	sTxs := querySettledTxs()
	for _, tx := range sTxs {
		if tx.OrderOut == 0 || tx.PriceOut != "" {
			continue
		}
		res, err := GetOrder(tx.Symbol, tx.OrderOut)
		if err != nil {
			fmt.Println("CheckOrd sell GetOrder err", err)
			continue
		}
		if res.Status != "FILLED" {
			continue
		}

		// 计算利息
		interest := calcInterest(tx.Quantity, res.Price, res.CummulativeQuoteQuantity)
		body := map[string]interface{}{
			"price_out":  res.Price,
			"interest":   interest,
			"settled_at": time.Now(),
		}
		err = UpdateTx(fmt.Sprintf("order_in=%d", tx.OrderIn), body)
		if err != nil {
			fmt.Println("CheckOrd sell UpdateTx err: ", err, tx)
			continue
		}
		fmt.Println("CheckOrd sell UpdateTx succeed", tx)
	}

	// 处理买入订单
	txs, _ := queryUnSettledTxs()
	for _, tx := range txs {
		if tx.OrderInStatus != 0 {
			continue
		}

		res, err := GetOrder(tx.Symbol, tx.OrderIn)
		if err != nil {
			fmt.Println("CheckOrd GetOrder err", err)
			continue
		}
		if res.Status != "FILLED" {
			continue
		}
		err = UpdateTx(fmt.Sprintf("order_in=%d", tx.OrderIn), map[string]interface{}{
			"order_in_status": 1, // 0委托，1成交，2取消
		})
		if err != nil {
			fmt.Println("checkOrd buy UpdateTx err: ", err, tx)
			continue
		}
		fmt.Println("checkOrd buy UpdateTx succeed", tx)
	}
}

func calcInterest(quantity, price, lastAmount string) string {
	q, _ := strconv.ParseFloat(quantity, 64)
	p, _ := strconv.ParseFloat(price, 64)
	la, _ := strconv.ParseFloat(lastAmount, 64)
	return fmt.Sprintf("%f", la-q*p-la*0.00075)
}

func checkSell(txs []*BnTxs, _smPri string) {
	shouldSell := false
	var tx *BnTxs
	for _, v := range txs {
		if v.PriceIn == _smPri && v.OrderInStatus == 1 && v.OrderOutStatus == 0 { // 价格匹配&&委托成功&&未建立委托
			shouldSell = true
			tx = v
			break
		}
	}
	if !shouldSell {
		return
	}
	ordRs, err := CreateOrder(tx.Symbol, tx.Quantity, _smPri, binance.SideTypeSell)
	if err != nil {
		fmt.Println("AutoGo CreateOrder sell err: ", err)
		return
	}
	err = UpdateTx(fmt.Sprintf("order_in=%d", tx.OrderIn), map[string]interface{}{
		"order_out":        ordRs.OrderID,
		"order_out_status": 2, // 1委托，2成交，3取消
		"status":           1, // 因为使用购买价卖出，因此能够马上卖出
	})
	if err != nil {
		fmt.Println("AutoGo UpdateTx sell err:", err, tx)
		return
	}
	fmt.Println("AutoGo CreateTx sell succeed", tx)
}

func checkBuy(txs []*BnTxs, smPri, symbol string) {
	shouldBuy := true
	for _, v := range txs {
		if v.PriceIn == smPri { // 如果存在小一点的价格，则不买了
			fmt.Println("checkBuy with order", v, smPri)
			shouldBuy = false
			break
		}
	}
	if !shouldBuy {
		return
	}
	quantity := gPriceToQuantityMap[smPri]
	ordRs, err := CreateOrder(symbol, quantity, smPri, binance.SideTypeBuy)
	if err != nil {
		fmt.Println("AutoGo CreateOrder buy err: ", err, quantity, smPri)
		return
	}
	tx := BnTxs{
		Symbol:   symbol,
		Quantity: quantity,
		PriceIn:  smPri,
		OrderIn:  ordRs.OrderID,
	}
	if err := CreateTx(&tx); err != nil {
		fmt.Println("AutoGo CreateTx buy err: ", err, tx)
		return
	}
	fmt.Println("AutoGo CreateTx buy succeed", tx)
}
