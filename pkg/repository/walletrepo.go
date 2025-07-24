package repository

import (
	"context"
	"time"

	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IWalletRepository interface {
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	CreateOrUpdateWallet(ctx context.Context, wallet *model.Wallet) error
	AddTransaction(ctx context.Context, tx *model.WalletTransaction) error
}

type WalletRepository struct {
	walletCol      *mongo.Collection
	transactionCol *mongo.Collection
}

func NewWalletRepository(walletCol, transactionCol *mongo.Collection) IWalletRepository {
	return &WalletRepository{walletCol: walletCol, transactionCol: transactionCol}
}

func (r *WalletRepository) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.walletCol.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) CreateOrUpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	_, err := r.walletCol.UpdateOne(ctx,
		bson.M{"user_id": wallet.UserID},
		bson.M{"$set": bson.M{
			"balance":    wallet.Balance,
			"updated_at": time.Now(),
		}},

		options.Update().SetUpsert(true),
	)
	return err
}

func (r *WalletRepository) AddTransaction(ctx context.Context, tx *model.WalletTransaction) error {
	tx.CreatedAt = time.Now()
	_, err := r.transactionCol.InsertOne(ctx, tx)
	return err
}
