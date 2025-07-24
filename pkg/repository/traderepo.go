package repository

import (
	"context"

	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITradeRepository interface {
	CreateTrade(trade *model.Trade) error
	GetTradeByID(id string) (*model.Trade, error)
	UpdateTrade(trade *model.Trade) error
	GetTradesByUserID(ctx context.Context, userID string) ([]*model.Trade, error)
}

type TradeRepository struct {
	db *mongo.Collection
}

func NewTradeRepository(col *mongo.Collection) ITradeRepository {
	return &TradeRepository{db: col}
}

func (r *TradeRepository) CreateTrade(trade *model.Trade) error {
	_, err := r.db.InsertOne(context.TODO(), trade)
	return err
}

func (r *TradeRepository) GetTradeByID(id string) (*model.Trade, error) {
	var trade model.Trade
	err := r.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&trade)
	if err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *TradeRepository) UpdateTrade(trade *model.Trade) error {
	_, err := r.db.ReplaceOne(context.TODO(), bson.M{"_id": trade.ID}, trade)
	return err
}

func (r *TradeRepository) GetTradesByUserID(ctx context.Context, userID string) ([]*model.Trade, error) {
	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"user_a_id": userID},
					{"user_b_id": userID},
				},
			},
			{
				"status": bson.M{"$nin": []string{"cancelled", "completed"}},
			},
		},
	}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trades []*model.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
