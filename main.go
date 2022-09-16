package main

import (
	"bn-auto/api"
	"time"
)

func main() {
	symbol := "ATOMUSDT"
	for {
		api.CheckOrd(symbol)
		api.AutoGo(symbol)
		time.Sleep(time.Second * 20)
	}
}
