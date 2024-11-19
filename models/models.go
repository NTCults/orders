package models

import (
	"errors"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	OofShard          string    `json:"oof_shard"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	Delivery          *Delivery `json:"delivery"`
	Payment           *Payment  `json:"payment"`
	Items             []*Item   `json:"items"`
}

func (o *Order) Validte() error {
	if o.Delivery == nil {
		return errors.New("delivery field is empty")
	}
	if o.Delivery == nil {
		return errors.New("payment field is empty")
	}
	if o.Items == nil {
		return errors.New("items field is empty")
	}
	return nil
}

type Delivery struct {
	Name    string `json:"order_uid"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Bank         string `json:"bank"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChartID     int    `json:"chart_id"`
	Price       int    `json:"price"`
	Sale        int    `json:"sale"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	TrackNumber string `json:"track_number"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	Brand       string `json:"brand"`
	Status      string `json:"status"`
}
