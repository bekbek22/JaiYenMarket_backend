package handler

import (
	"github.com/bekbek22/JaiYenMarket_backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type IWalletHandler interface {
	GetBalance(c *fiber.Ctx) error
	Deposit(c *fiber.Ctx) error
	Withdraw(c *fiber.Ctx) error
	Transfer(c *fiber.Ctx) error
}

type WalletHandler struct {
	service service.IWalletService
}

func NewWalletHandler(s service.IWalletService) IWalletHandler {
	return &WalletHandler{service: s}
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	balance, err := h.service.GetBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get balance"})
	}
	return c.JSON(fiber.Map{"balance": balance})
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

func (h *WalletHandler) Deposit(c *fiber.Ctx) error {
	var req DepositRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	userID := c.Locals("user_id").(string)
	err := h.service.Deposit(c.Context(), userID, req.Amount, req.Note)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deposit successful"})
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

func (h *WalletHandler) Withdraw(c *fiber.Ctx) error {
	var req WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	userID := c.Locals("user_id").(string)
	err := h.service.Withdraw(c.Context(), userID, req.Amount, req.Note)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "withdraw requested"})
}

type TransferRequest struct {
	ToUserID   string  `json:"to_user_id"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
	FeePercent float64 `json:"fee_percent"` // เช่น 5.0 หมายถึง 5%
}

func (h *WalletHandler) Transfer(c *fiber.Ctx) error {
	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	fromUserID := c.Locals("user_id").(string)
	err := h.service.Transfer(c.Context(), fromUserID, req.ToUserID, req.Amount, req.Note, req.FeePercent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "transfer successful"})
}
