package services

import (
	"errors"
	"time"

	"stream-app/pkg/models"
	"stream-app/pkg/utils"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	secretKey string
}

func NewAuthService(db *gorm.DB, secretKey string) *AuthService {
	return &AuthService{
		db:        db,
		secretKey: secretKey,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("пользователь не найден")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	token, err := utils.GenerateJWT(user.ID, s.secretKey, time.Hour*72)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {

	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:       utils.GenerateUUID(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
