package handlers

import (
	"net/http"

	"stream-app/pkg/services"

	"github.com/gin-gonic/gin"
)

type StreamHandler struct {
	service services.StreamService
}

func NewStreamHandler(service services.StreamService) *StreamHandler {
	return &StreamHandler{service: service}
}

// CreateStreamHandler обрабатывает запросы на создание нового потока
func (h *StreamHandler) CreateStreamHandler(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
		return
	}

	// Получаем ID текущего пользователя из контекста (предполагается, что пользователь аутентифицирован)
	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима аутентификация"})
		return
	}

	stream, err := h.service.CreateStream(req.Title, authorID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stream)
}

// GetStreamHandler обрабатывает запросы на получение информации о конкретном потоке
func (h *StreamHandler) GetStreamHandler(c *gin.Context) {
	id := c.Param("id")

	stream, err := h.service.GetStreamByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Поток не найден"})
		return
	}

	c.JSON(http.StatusOK, stream)
}

// GetAllStreamsHandler обрабатывает запросы на получение списка всех потоков
func (h *StreamHandler) GetAllStreamsHandler(c *gin.Context) {
	streams, err := h.service.GetAllStreams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список потоков"})
		return
	}

	c.JSON(http.StatusOK, streams)
}

// DeleteStreamHandler обрабатывает запросы на удаление потока
func (h *StreamHandler) DeleteStreamHandler(c *gin.Context) {
	id := c.Param("id")

	// Получаем ID текущего пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходима аутентификация"})
		return
	}

	// Получаем поток для проверки принадлежности
	stream, err := h.service.GetStreamByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Поток не найден"})
		return
	}

	if stream.AuthorID.String() != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на удаление этого потока"})
		return
	}

	if err := h.service.DeleteStream(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить поток"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Поток успешно удалён"})
}
