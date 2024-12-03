package entity

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required" db:"order_uid"`
	TrackNumber       string    `json:"track_number" validate:"required" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery" validate:"required" db:"-"`
	Payment           Payment   `json:"payment" validate:"required" db:"-"`
	Items             []Item    `json:"items" validate:"required" db:"-"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	ShardKey          string    `json:"shardkey" db:"shardkey"`
	SMID              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OOFShard          string    `json:"oof_shard" db:"oof_shard"`
}
