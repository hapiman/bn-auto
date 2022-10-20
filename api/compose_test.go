package api

import (
	"fmt"
	"strconv"
	"testing"
)

const gridAmount = 20

func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func TestAutoGo(t *testing.T) {
	for i := 7.0; i < 8; i += 0.1 {
		fmt.Printf("\"%v\":\"%v\",\n", decimal(i), decimal(gridAmount/i))
		// fmt.Printf("\"%v\",", decimal(i))
	}
	for i := float64(18 + 0.25); i <= 20; i += 0.25 {
		fmt.Printf("\"%v\":\"%v\",\n", decimal(i), decimal(gridAmount/i))
	}
}

func Test_calcInterest(t *testing.T) {
	type args struct {
		quantity string
		price    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test_calcInterest1", args{
			quantity: "1.28",
			price:    "1.56",
		}, "1.996800,"},
		{"Test_calcInterest1", args{
			quantity: "1.28",
			price:    "1.58002",
		}, "2.022426,"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcInterest(tt.args.quantity, tt.args.price, ""); got != tt.want {
				t.Errorf("calcInterest() = %v, want %v", got, tt.want)
			}
		})
	}
}
