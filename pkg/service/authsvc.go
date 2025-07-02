package service

import (
	"errors"

	"github.com/bekbek22/JaiYenMarket_backend/config"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/repository"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/utils"
)

type IAuthService interface {
	Register(user *model.User) error
	Login(email, password string) (string, error)
}

type AuthService struct {
	repo repository.IAuthRepository
	Cfg  *config.Config
}

func NewAuthService(r repository.IAuthRepository, cfg *config.Config) IAuthService {
	return &AuthService{
		repo: r,
		Cfg:  cfg,
	}
}

func (s *AuthService) Register(user *model.User) error {
	exisiting, _ := s.repo.FindByEmail(user.Email)

	if exisiting != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	Users := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Role:     "user", //set default
		Password: hashedPassword,
	}

	return s.repo.CreateUser(Users)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role, s.Cfg.JWTSecret)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
