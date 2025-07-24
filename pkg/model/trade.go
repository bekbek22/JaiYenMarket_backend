package model

import "time"

type TradeStatus string

const (
	TradePending   TradeStatus = "pending"
	TradeConfirmed TradeStatus = "confirmed"
	TradeCompleted TradeStatus = "completed"
	TradeCancelled TradeStatus = "cancelled"
)

type Item struct {
	Type string `bson:"type" json:"type"`
	Name string `bson:"name" json:"name"`
	Code string `bson:"code" json:"code"`
	Note string `bson:"note,omitempty" json:"note,omitempty"`
}

type Trade struct {
	ID        string      `bson:"_id" json:"id"`
	UserAID   string      `bson:"user_a_id" json:"user_a_id"`
	UserBID   string      `bson:"user_b_id" json:"user_b_id"`
	ItemFromA []Item      `bson:"item_from_a" json:"item_from_a"`
	ItemFromB []Item      `bson:"item_from_b" json:"item_from_b"`
	ConfirmA  bool        `bson:"confirm_a" json:"confirm_a"`
	ConfirmB  bool        `bson:"confirm_b" json:"confirm_b"`
	Status    TradeStatus `bson:"status" json:"status"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}
