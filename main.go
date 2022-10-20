package main

import (
	"bn-auto/api"
	"fmt"
	"time"
)

func main() {
	symbol := "ATOMBUSD"
	symbol2 := "APTBUSD"
	for {
		api.CheckOrd(symbol)
		api.AutoGo(symbol)
		aptTrade(symbol2)
		time.Sleep(time.Second * 20)
	}
}

func aptTrade(symbol string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("aptTrade err: ", err)
		}
	}()
	api.CheckOrd(symbol)
	api.AutoGo(symbol)

}
