package entity

import "encoding/json"

type Order struct {
	OrderUID          string `json:"order_uid" fake:"{uuid}"`
	TrackNumber       string `json:"track_number" fake:"WB{lettern:12}"`
	Entry             string `json:"entry" fake:"WB??"`
	Delivery          Delivery
	Payment           Payment
	Items             []Item `fakesize:"1,5"`
	Locale            string `json:"locale" fake:"{languagebcp}"`
	InternalSignature string `json:"internal_signature" fake:"skip"`
	CustomerID        string `json:"customer_id" fake:"test0###"`
	DeliveryService   string `json:"delivery_service" fake:"?????"`
	Shardkey          string `json:"shardkey" fake:"#"`
	SmID              int    `json:"sm_id" fake:"##"`
	DateCreated       string `json:"date_created" fake:"{pastdate}"`
	OofShard          string `json:"oof_shard" fake:"#"`
}

type Delivery struct {
	Name    string `json:"name" fake:"{name}"`
	Phone   string `json:"phone" fake:"+{phone}"`
	Zip     string `json:"zip" fake:"{zip}"`
	City    string `json:"city" fake:"{city}"`
	Address string `json:"address" fake:"{street}"`
	Region  string `json:"region" fake:"{state}"`
	Email   string `json:"email" fake:"{email}"`
}

type Payment struct {
	Transaction  string `json:"transaction" fake:"{uuid}"`
	RequestID    string `json:"request_id" fake:"skip"`
	Currency     string `json:"currency" fake:"{currencyshort}"`
	Provider     string `json:"provider" fake:"{creditcardtype}"`
	Amount       int    `json:"amount" fake:"{number:1,1000}"`
	PaymentDT    int    `json:"payment_dt" fake:"{digitn:10}"`
	Bank         string `json:"bank" fake:"{company:name}"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:20,300}"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:1,100}"`
	CustomFee    int    `json:"custom_fee" fake:"{number:1,30}"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" fake:"#######"`
	TrackNumber string `json:"track_number" fake:"{lettern:14}"`
	Price       int    `json:"price" fake:"{number:70,1000}"`
	RID         string `json:"rid" fake:"{lettern:21}"`
	Name        string `json:"name" fake:"{firstname}"`
	Sale        int    `json:"sale" fake:"{number:1,100}"`
	Size        string `json:"size" fake:"{randomstring:[XS,S,M,L,XL,XXL]}"`
	TotalPrice  int    `json:"total_price" fake:"{number:70,1000}"`
	NMID        int    `json:"nm_id" fake:"{digitn:8}"`
	Brand       string `json:"brand" fake:"{company}"`
	Status      int    `json:"status" fake:"{digitn:3}"`
}

func (o Order) String() string {
	jsonData, _ := json.MarshalIndent(o, "", "  ")
	return string(jsonData)
}
