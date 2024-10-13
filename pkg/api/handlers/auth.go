package handlers

import (
	"net/http"
	"stream-app/pkg/config"
	"stream-app/pkg/logger"
	"stream-app/pkg/services"
	"stream-app/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db      *gorm.DB
	cfg     *config.Config
	logger  *logger.Logger
	service *services.AuthService
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config, logger *logger.Logger) *AuthHandler {
	service := services.NewAuthService(db, cfg.JWT.SecretKey)
	return &AuthHandler{
		db:      db,
		cfg:     cfg,
		logger:  logger,
		service: service,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.SugaredLogger.Warnf("request error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.logger.SugaredLogger.Errorf("cannot login: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.SugaredLogger.Warnf("request error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.SugaredLogger.Errorf("Ошибка регистрации: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, h.cfg.JWT.SecretKey, time.Hour*72)
	if err != nil {
		h.logger.SugaredLogger.Errorf("Ошибка генерации токена: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
