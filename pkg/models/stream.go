package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stream struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	AuthorID  uuid.UUID      `gorm:"type:uuid;not null" json:"author_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Stream) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
