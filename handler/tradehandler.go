package handler

import (
	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type ITradeHandler interface {
	CreateTrade(c *fiber.Ctx) error
	AddItemToTrade(c *fiber.Ctx) error
	ConfirmTrade(c *fiber.Ctx) error
	GetTrade(c *fiber.Ctx) error
	UnconfirmTrade(c *fiber.Ctx) error
	CancelTrade(c *fiber.Ctx) error
	GetMyTrades(c *fiber.Ctx) error
}

type TradeHandler struct {
	service service.ITradeService
}

func NewTradeHandler(s service.ITradeService) ITradeHandler {
	return &TradeHandler{service: s}
}

// POST /api/trade
func (h *TradeHandler) CreateTrade(c *fiber.Ctx) error {
	var req struct {
		UserBID string `json:"user_b_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	userAID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	trade := &model.Trade{
		UserAID: userAID,
		UserBID: req.UserBID,
	}

	if err := h.service.CreateTrade(c.Context(), trade); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create trade"})
	}

	return c.JSON(fiber.Map{"message": "trade created", "trade_id": trade.ID})
}

func (h *TradeHandler) AddItemToTrade(c *fiber.Ctx) error {
	tradeID := c.Params("id")

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var req struct {
		Items []model.Item `json:"items"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "invalid input",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no items provided"})
	}

	err := h.service.AddItemToTrade(c.Context(), tradeID, userID, req.Items)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "items added to trade"})
}

// POST /api/trade/:id/confirm
func (h *TradeHandler) ConfirmTrade(c *fiber.Ctx) error {
	tradeID := c.Params("id")

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := h.service.ConfirmTrade(c.Context(), tradeID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "trade confirmation updated"})
}

func (h *TradeHandler) UnconfirmTrade(c *fiber.Ctx) error {
	tradeID := c.Params("id")

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := h.service.UnconfirmTrade(c.Context(), tradeID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "trade unconfirmed"})
}

func (h *TradeHandler) CancelTrade(c *fiber.Ctx) error {
	tradeID := c.Params("id")

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := h.service.CancelTrade(c.Context(), tradeID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "trade cancelled"})
}

// GET /api/trade/:id
func (h *TradeHandler) GetTrade(c *fiber.Ctx) error {
	tradeID := c.Params("id")

	trade, err := h.service.GetTrade(c.Context(), tradeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "trade not found"})
	}

	return c.JSON(trade)
}

func (h *TradeHandler) GetMyTrades(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	trades, err := h.service.GetTradesByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch trades"})
	}

	return c.JSON(trades)
}
