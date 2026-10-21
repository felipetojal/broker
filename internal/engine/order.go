package engine

import "time"

type OrderSide int8
const (
		Buy OrderSide = 1
		Sell OrderSide = 2
)

type Order struct {
		Id string
		OrderSide `json:"type"`
		Price int64 `json:"price"`
		Quantity int64 `json:"quantity"`
		In int64
		CreateAt time.Time
		Index int
		Status bool
		IdUser int64 `json:"id_user"`
		Ticker string `json:"ticker"`
}


type Trade struct {
		IdBuy string `json:"id_buy"`
		IdSell string `json:"id_sell"`
		Price int64 `json:"price"`
		Quantity int64 `json:"quantity"`
		TradeAt time.Time
}
