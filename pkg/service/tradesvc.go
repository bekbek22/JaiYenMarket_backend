package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/repository"
	"github.com/google/uuid"
)

type ITradeService interface {
	CreateTrade(ctx context.Context, trade *model.Trade) error
	ConfirmTrade(ctx context.Context, tradeID, userID string) error
	UnconfirmTrade(ctx context.Context, tradeID, userID string) error
	CancelTrade(ctx context.Context, tradeID, userID string) error
	GetTrade(ctx context.Context, tradeID string) (*model.Trade, error)
	AddItemToTrade(ctx context.Context, tradeID, userID string, items []model.Item) error
	GetTradesByUserID(ctx context.Context, userID string) ([]*model.Trade, error)
}

type TradeService struct {
	repo repository.ITradeRepository
}

func NewTradeService(r repository.ITradeRepository) ITradeService {
	return &TradeService{repo: r}
}

func (s *TradeService) CreateTrade(ctx context.Context, trade *model.Trade) error {
	trade.ID = uuid.NewString()
	trade.Status = "pending"
	trade.ConfirmA = false
	trade.ConfirmB = false
	trade.CreatedAt = time.Now()
	trade.UpdatedAt = time.Now()

	return s.repo.CreateTrade(trade)
}

func (s *TradeService) AddItemToTrade(ctx context.Context, tradeID, userID string, items []model.Item) error {
	trade, err := s.repo.GetTradeByID(tradeID)
	if err != nil {
		return fmt.Errorf("trade not found: %w", err)
	}

	if trade.Status != "pending" {
		return errors.New("cannot add items to a non-pending trade")
	}

	if (trade.UserAID == userID && trade.ConfirmA) || (trade.UserBID == userID && trade.ConfirmB) {
		return errors.New("you already confirmed, cannot add items")
	}

	//ถ้า A เพิ่มของ
	if trade.UserAID == userID {
		trade.ItemFromA = append(trade.ItemFromA, items...)
	} else if trade.UserBID == userID {
		trade.ItemFromB = append(trade.ItemFromB, items...)
	} else {
		return errors.New("user not part of this trade")
	}

	trade.UpdatedAt = time.Now()

	return s.repo.UpdateTrade(trade)
}

func (s *TradeService) ConfirmTrade(ctx context.Context, tradeID, userID string) error {
	trade, err := s.repo.GetTradeByID(tradeID)
	if err != nil {
		return err
	}

	if trade.Status != "pending" {
		return errors.New("trade not pending")
	}

	updated := false
	if trade.UserAID == userID {
		trade.ConfirmA = true
		updated = true
	} else if trade.UserBID == userID {
		trade.ConfirmB = true
		updated = true
	}

	if !updated {
		return errors.New("user not in this trade")
	}

	if trade.ConfirmA && trade.ConfirmB {
		trade.Status = "completed"
	}
	trade.UpdatedAt = time.Now()
	return s.repo.UpdateTrade(trade)
}

func (s *TradeService) UnconfirmTrade(ctx context.Context, tradeID, userID string) error {
	trade, err := s.repo.GetTradeByID(tradeID)
	if err != nil {
		return fmt.Errorf("trade not found: %w", err)
	}

	if trade.Status != "pending" {
		return errors.New("cannot unconfirm a non-pending trade")
	}

	if trade.UserAID == userID {
		trade.ConfirmA = false
	} else if trade.UserBID == userID {
		trade.ConfirmB = false
	} else {
		return errors.New("user not part of this trade")
	}

	trade.UpdatedAt = time.Now()
	return s.repo.UpdateTrade(trade)
}

func (s *TradeService) CancelTrade(ctx context.Context, tradeID, userID string) error {
	trade, err := s.repo.GetTradeByID(tradeID)
	if err != nil {
		return fmt.Errorf("trade not found: %w", err)
	}

	// ตรวจสอบว่า user เป็นเจ้าของ trade ฝั่งใดฝั่งหนึ่ง
	if trade.UserAID != userID && trade.UserBID != userID {
		return errors.New("you are not part of this trade")
	}

	// ตรวจสอบว่าที่ยกเลิกได้ต้องยังไม่ complete หรือ cancel ไปแล้ว
	if trade.Status != "pending" {
		return errors.New("cannot cancel a trade that is not pending")
	}

	trade.Status = "cancelled"
	trade.UpdatedAt = time.Now()

	return s.repo.UpdateTrade(trade)
}

func (s *TradeService) GetTrade(ctx context.Context, tradeID string) (*model.Trade, error) {
	return s.repo.GetTradeByID(tradeID)
}

func (s *TradeService) GetTradesByUserID(ctx context.Context, userID string) ([]*model.Trade, error) {
	return s.repo.GetTradesByUserID(ctx, userID)
}
