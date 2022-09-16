package api

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"os"
)

var client *binance.Client

func init() {
	var apiKey = os.Getenv("BnApiKey")
	var secretKey = os.Getenv("BnSecretKey")
	client = binance.NewClient(apiKey, secretKey)
}

func CreateOrder(symbol, quantity, price string, side binance.SideType) (res *binance.CreateOrderResponse, err error) {
	return client.NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
}

func ListOrders(symbol string) (res []*binance.Order, err error) { // 需要挂起来订单才能看到
	return client.NewListOpenOrdersService().Symbol(symbol).
		Do(context.Background())
}

func ListOpenOrders(symbol string) (res []*binance.Order, err error) { // 需要挂起来的订单才能看到
	return client.NewListOpenOrdersService().Symbol(symbol).
		Do(context.Background())
}

func ListTickerPrice(symbols []string) (pMap map[string]string) {
	pMap = map[string]string{}
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	sMap := map[string]struct{}{}
	for _, sb := range symbols {
		sMap[sb] = struct{}{}
	}
	for _, p := range prices {
		if _, ok := sMap[p.Symbol]; ok {
			pMap[p.Symbol] = p.Price
		}
	}
	return
}

func GetPrice(symbol string) string {
	mp := ListTickerPrice([]string{symbol})
	_, ok := mp[symbol]
	if !ok {
		panic("unexpected error")
	}
	return mp[symbol]
}
