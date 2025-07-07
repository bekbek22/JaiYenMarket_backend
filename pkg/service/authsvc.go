package service

import (
	"errors"
	"time"

	"github.com/bekbek22/JaiYenMarket_backend/config"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/model"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/repository"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthService interface {
	Register(user *model.User) error
	Login(email, password string) (string, int64, *model.User, error)
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

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("กรุณากรอกช้อมูลให้ครบ")
	}

	exisiting, _ := s.repo.FindByEmail(user.Email)
	if exisiting != nil {
		return errors.New("อีเมลนี้ถูกใช้งานแล้ว")
	}

	if user.Role == "" {
		user.Role = "user"
	}

	// สร้าง ObjectID ใหม่
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().Unix()

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	Users := &model.User{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Password: hashedPassword,
	}

	err = s.repo.CreateUser(Users)
	if err != nil {
		return errors.New("การบันทึกข้อมูลไม่สำเร็จ")
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, int64, *model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", 0, nil, errors.New("ไม่มีผู้ใช้งานนี้")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", 0, nil, errors.New("รหัสผ่านไม่ถูกต้อง")
	}
	expiresIn := int64(3600)
	exp := time.Now().Add(time.Second * time.Duration(expiresIn)).Unix() // 3600 sec

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role, s.Cfg.JWTSecret, exp)
	if err != nil {
		return "", 0, nil, errors.New("failed to generate token")
	}

	return token, expiresIn, user, nil
}
