package model

import "time"

type Wallet struct {
	UserID    string    `bson:"user_id"`
	Balance   float64   `bson:"balance"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type WalletTransaction struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	Amount    float64   `bson:"amount"`
	Type      string    `bson:"type"`
	Note      string    `bson:"note"`
	CreatedAt time.Time `bson:"created_at"`
}
