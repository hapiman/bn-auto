## 依赖模块

- [go-binance](https://github.com/adshao/go-binance)

## 数据库 
存放订单数据，便于记录买卖操作和计算利息 
```sql
CREATE TABLE `bn_txs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `symbol` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易对',
  `quantity` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '交易数量',
  `price_in` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '买单价格',
  `price_out` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '卖单价格',
  `order_in` int(255) DEFAULT NULL COMMENT '买单订单号',
  `order_in_status` int(255) DEFAULT '0' COMMENT '0委托，1成交，2取消',
  `order_out` int(255) DEFAULT NULL COMMENT '卖单订单号',
  `order_out_status` int(255) DEFAULT '0' COMMENT '1委托，2成交，3取消',
  `status` int(11) DEFAULT '0' COMMENT '结束状态，0未结束，1已结束',
  `interest` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '利润',
  `created_at` datetime DEFAULT NULL COMMENT '买单委托时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
```