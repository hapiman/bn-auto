package api

import (
	"fmt"
	"github.com/adshao/go-binance/v2"
)

// 交易网格
var gPriceToQuantityMap = map[string]string{
	"10":   "2",
	"10.2": "1.96",
	"10.4": "1.92",
	"10.6": "1.89",
	"10.8": "1.85",
	"11":   "1.82",
	"11.2": "1.79",
	"11.4": "1.75",
	"11.6": "1.72",
	"11.8": "1.69",
	"12":   "1.67",
	"12.2": "1.64",
	"12.4": "1.61",
	"12.6": "1.59",
	"12.8": "1.56",
	"13":   "1.54",
	"13.2": "1.52",
	"13.4": "1.49",
	"13.6": "1.47",
	"13.8": "1.45",
	"14":   "1.43",
	"14.2": "1.41",
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

func queryUnSettledTxs() (txs []*BnTxs) {
	err := gDb.Debug().Raw("select * from bn_txs where status=0").Scan(&txs).Error
	if err != nil {
		fmt.Println("queryUnSettledTxs err", err)
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
func AutoGo() {
	symbol := "ATOMUSDT"
	price := GetPrice(symbol)
	smPri, _ := getLeftAndRightPrice(price)  // 14.8 => 14.6
	_smPri, _ := getLeftAndRightPrice(smPri) // 14.6  => 14.4
	txs := queryUnSettledTxs()
	fmt.Println("AutoGo unsettled order length: ", len(txs))
	for _, v := range txs {
		if v.PriceIn != smPri { // 没有未结束的买单
			quantity := gPriceToQuantityMap[smPri]
			ordRs, err := CreateOrder(symbol, quantity, smPri, binance.SideTypeBuy)
			if err != nil {
				fmt.Println("AutoGo CreateOrder buy err: ", err, v)
				break
			}
			tx := BnTxs{
				Symbol:        symbol,
				Quantity:      quantity,
				PriceIn:       smPri,
				OrderIn:       ordRs.OrderID,
				OrderInStatus: 1,
			}
			if err := CreateTx(&tx); err != nil {
				fmt.Println("AutoGo CreateTx buy err: ", err, v)
			}
			break
		}

		if v.PriceIn == _smPri { // 更新操作
			ordRs, err := CreateOrder(symbol, v.Quantity, smPri, binance.SideTypeSell)
			if err != nil {
				fmt.Println("AutoGo CreateOrder sell err: ", err, v)
				break
			}
			err = UpdateTx(fmt.Sprintf("order_in=%d", v.OrderIn), map[string]interface{}{
				"order_out":        ordRs.OrderID,
				"order_out_status": 1,
				"status":           1,
			})
			if err != nil {
				fmt.Println("AutoGo UpdateTx sell err:", err, v)
				break
			}
		}
	}
}
