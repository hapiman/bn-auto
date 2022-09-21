package main

import (
	"bn-auto/api"
	"time"
)

func main() {
	symbol := "ATOMBUSD"
	for {
		api.CheckOrd()
		api.AutoGo(symbol)
		time.Sleep(time.Second * 20)
	}
}
