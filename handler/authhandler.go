package handler

import (
	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type AuthHandler struct {
	service service.IAuthService
}

func NewAuthHandler(s service.IAuthService) IAuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	if err := h.service.Register(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	response := model.RegisterResponse{
		Status:  fiber.StatusCreated,
		Message: "สร้างผู้ใช้สำเร็จ",
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	return c.Status(response.Status).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var creds model.User
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	token, expiresIn, user, err := h.service.Login(creds.Email, creds.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	response := model.LoginResponse{
		Status:      fiber.StatusOK,
		Message:     "เข้าสู่ระบบสำเร็จ",
		AccessToken: token,
		ExpiresIn:   int(expiresIn),
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	return c.Status(response.Status).JSON(response)
}
