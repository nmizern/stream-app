package services

import (
	"errors"

	"stream-app/pkg/models"
	"stream-app/pkg/repositories"

	"github.com/google/uuid"
)

type StreamService interface {
	CreateStream(title string, authorID string) (*models.Stream, error)
	GetStreamByID(id string) (*models.Stream, error)
	GetAllStreams() ([]models.Stream, error)
	DeleteStream(id string) error
}

type streamService struct {
	repo repositories.StreamRepository
}

func NewStreamService(repo repositories.StreamRepository) StreamService {
	return &streamService{repo: repo}
}

func (s *streamService) CreateStream(title string, authorID string) (*models.Stream, error) {

	if title == "" {
		return nil, errors.New("название потока не может быть пустым")
	}

	if _, err := uuid.Parse(authorID); err != nil {
		return nil, errors.New("некорректный ID автора")
	}

	stream := &models.Stream{
		Title:    title,
		AuthorID: uuid.MustParse(authorID),
	}

	if err := s.repo.Create(stream); err != nil {
		return nil, err
	}

	return stream, nil
}

func (s *streamService) GetStreamByID(id string) (*models.Stream, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("некорректный ID потока")
	}

	stream, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (s *streamService) GetAllStreams() ([]models.Stream, error) {
	return s.repo.GetAll()
}

func (s *streamService) DeleteStream(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("некорректный ID потока")
	}

	return s.repo.Delete(id)
}
