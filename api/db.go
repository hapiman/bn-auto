package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var gDb *gorm.DB

func init() {
	dsn := "root:abc123@tcp(127.0.0.1:3306)/cosmosdb?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = gDsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gDb = db
}

// BnTxs 币安交易记录
type BnTxs struct {
	Id             int64     `json:"id"`
	Symbol         string    `json:"symbol"`
	Quantity       string    `json:"quantity"`
	PriceIn        string    `json:"price_in"`
	PriceOut       string    `json:"price_out"`
	OrderIn        int64     `json:"order_in"`
	OrderInStatus  int64     `json:"order_in_status"`
	OrderOut       int64     `json:"order_out"`
	OrderOutStatus int64     `json:"order_out_status"`
	status         int64     `json:"status"`
	interest       string    `json:"interest"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 创建
func CreateTx(tx *BnTxs) error {
	return gDb.Create(tx).Error
}

// 更新
func UpdateTx(cons string, attrs map[string]interface{}) error {
	return gDb.Table("bn_txs").Where(cons).Updates(attrs).Error
}
