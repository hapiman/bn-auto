package api

import (
	"fmt"
	"strconv"
	"testing"
)

func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func TestAutoGo(t *testing.T) {
	for i := 10.0; i <= 18; i += 0.2 {
		// fmt.Printf("\"%v\":\"%v\",\n", decimal(i), decimal(20/i))
		fmt.Printf("\"%v\",", decimal(i))
	}
}
