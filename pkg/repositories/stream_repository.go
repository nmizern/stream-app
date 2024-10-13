package repositories

import (
	"stream-app/pkg/models"

	"gorm.io/gorm"
)

type StreamRepository interface {
	Create(stream *models.Stream) error
	GetByID(id string) (*models.Stream, error)
	GetAll() ([]models.Stream, error)
	Delete(id string) error
}

type streamRepository struct {
	db *gorm.DB
}

func NewStreamRepository(db *gorm.DB) StreamRepository {
	return &streamRepository{db: db}
}

func (r *streamRepository) Create(stream *models.Stream) error {
	return r.db.Create(stream).Error
}

func (r *streamRepository) GetByID(id string) (*models.Stream, error) {
	var stream models.Stream
	if err := r.db.First(&stream, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &stream, nil
}

func (r *streamRepository) GetAll() ([]models.Stream, error) {
	var streams []models.Stream
	if err := r.db.Find(&streams).Error; err != nil {
		return nil, err
	}
	return streams, nil
}

func (r *streamRepository) Delete(id string) error {
	return r.db.Delete(&models.Stream{}, "id = ?", id).Error
}
