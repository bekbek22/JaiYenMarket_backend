package service

import (
	"context"
	"errors"
	"time"

	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/repository"
)

type IWalletService interface {
	GetBalance(ctx context.Context, userID string) (float64, error)
	Deposit(ctx context.Context, userID string, amount float64, note string) error
	Withdraw(ctx context.Context, userID string, amount float64, note string) error
	Transfer(ctx context.Context, fromUserID, toUserID string, amount float64, note string, feePercent float64) error
}

type WalletService struct {
	repo repository.IWalletRepository
}

func NewWalletService(repo repository.IWalletRepository) IWalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(ctx context.Context, userID string) (float64, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *WalletService) Deposit(ctx context.Context, userID string, amount float64, note string) error {
	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}
	wallet, _ := s.repo.GetWalletByUserID(ctx, userID)
	if wallet == nil {
		wallet = &model.Wallet{UserID: userID, Balance: 0, UpdatedAt: time.Now()}
	}
	wallet.Balance += amount
	err := s.repo.CreateOrUpdateWallet(ctx, wallet)
	if err != nil {
		return err
	}
	return s.repo.AddTransaction(ctx, &model.WalletTransaction{
		UserID: userID,
		Amount: amount,
		Type:   "deposit",
		Note:   note,
	})
}

func (s *WalletService) Withdraw(ctx context.Context, userID string, amount float64, note string) error {
	if amount <= 0 {
		return errors.New("invalid withdraw amount")
	}
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil || wallet.Balance < amount {
		return errors.New("insufficient balance")
	}
	wallet.Balance -= amount
	if err := s.repo.CreateOrUpdateWallet(ctx, wallet); err != nil {
		return err
	}
	return s.repo.AddTransaction(ctx, &model.WalletTransaction{
		UserID: userID,
		Amount: -amount,
		Type:   "withdraw",
		Note:   note,
	})
}

func (s *WalletService) Transfer(ctx context.Context, fromUserID, toUserID string, amount float64, note string, feePercent float64) error {
	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}
	walletFrom, err := s.repo.GetWalletByUserID(ctx, fromUserID)
	if err != nil || walletFrom.Balance < amount {
		return errors.New("insufficient balance")
	}
	fee := amount * feePercent / 100
	toReceive := amount - fee

	walletTo, _ := s.repo.GetWalletByUserID(ctx, toUserID)
	if walletTo == nil {
		walletTo = &model.Wallet{UserID: toUserID, Balance: 0}
	}

	// Update balances
	walletFrom.Balance -= amount
	walletTo.Balance += toReceive

	if err := s.repo.CreateOrUpdateWallet(ctx, walletFrom); err != nil {
		return err
	}
	if err := s.repo.CreateOrUpdateWallet(ctx, walletTo); err != nil {
		return err
	}

	// Transactions
	_ = s.repo.AddTransaction(ctx, &model.WalletTransaction{
		UserID: fromUserID,
		Amount: -amount,
		Type:   "trade_payment",
		Note:   note,
	})

	_ = s.repo.AddTransaction(ctx, &model.WalletTransaction{
		UserID: toUserID,
		Amount: toReceive,
		Type:   "trade_receive",
		Note:   note,
	})

	if fee > 0 {
		_ = s.repo.AddTransaction(ctx, &model.WalletTransaction{
			UserID: "system",
			Amount: fee,
			Type:   "fee",
			Note:   "platform fee",
		})
	}

	return nil
}
