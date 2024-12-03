package entities

import "encoding/json"

type Order struct {
	OrderUID          string `json:"order_uid" fake:"{uuid}" validate:"required,uuid"`
	TrackNumber       string `json:"track_number" fake:"WB{lettern:12}" validate:"required,alphanum,len=14"`
	Entry             string `json:"entry" fake:"WB??" validate:"required,alphanum,len=4"`
	Delivery          Delivery
	Payment           Payment
	Items             []Item `fakesize:"1,5"`
	Locale            string `json:"locale" fake:"{languagebcp}" validate:"required"`
	InternalSignature string `json:"internal_signature" fake:"skip"`
	CustomerID        string `json:"customer_id" fake:"test0###" validate:"required,alphanum"`
	DeliveryService   string `json:"delivery_service" fake:"?????" validate:"required,alphanum"`
	Shardkey          string `json:"shardkey" fake:"#" validate:"required,numeric"`
	SmID              int    `json:"sm_id" fake:"##" validate:"required,numeric"`
	DateCreated       string `json:"date_created" fake:"{pastdate}" validate:"required"`
	OofShard          string `json:"oof_shard" fake:"#" validate:"required,numeric"`
}

type Delivery struct {
	Name    string `json:"name" fake:"{name}" validate:"required"`
	Phone   string `json:"phone" fake:"+{phone}" validate:"required"`
	Zip     string `json:"zip" fake:"{zip}" validate:"required,numeric"`
	City    string `json:"city" fake:"{city}" validate:"required"`
	Address string `json:"address" fake:"{street}" validate:"required"`
	Region  string `json:"region" fake:"{state}" validate:"required"`
	Email   string `json:"email" fake:"{email}" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" fake:"{uuid}" validate:"required,uuid"`
	RequestID    string `json:"request_id" fake:"skip"`
	Currency     string `json:"currency" fake:"{currencyshort}" validate:"required"`
	Provider     string `json:"provider" fake:"{creditcardtype}" validate:"required"`
	Amount       int    `json:"amount" fake:"{number:1,1000}" validate:"required,numeric"`
	PaymentDT    int    `json:"payment_dt" fake:"{digitn:10}" validate:"required,numeric"`
	Bank         string `json:"bank" fake:"{company:name}" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:20,300}" validate:"required,numeric"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:1,100}" validate:"required,numeric"`
	CustomFee    int    `json:"custom_fee" fake:"{number:1,30}" validate:"required,numeric"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" fake:"#######" validate:"required,numeric"`
	TrackNumber string `json:"track_number" fake:"{lettern:14}" validate:"required,alphanum,len=14"`
	Price       int    `json:"price" fake:"{number:70,1000}" validate:"required,numeric"`
	RID         string `json:"rid" fake:"{lettern:21}" validate:"required,alphanum,len=21"`
	Name        string `json:"name" fake:"{firstname}" validate:"required"`
	Sale        int    `json:"sale" fake:"{number:1,100}" validate:"required,numeric"`
	Size        string `json:"size" fake:"{randomstring:[XS,S,M,L,XL,XXL]}" validate:"required"`
	TotalPrice  int    `json:"total_price" fake:"{number:70,1000}" validate:"required,numeric"`
	NMID        int    `json:"nm_id" fake:"{digitn:8}" validate:"required,numeric"`
	Brand       string `json:"brand" fake:"{company}" validate:"required"`
	Status      int    `json:"status" fake:"{digitn:3}" validate:"required,numeric"`
}

func (o Order) String() string {
	jsonData, _ := json.MarshalIndent(o, "", "  ")
	return string(jsonData)
}
