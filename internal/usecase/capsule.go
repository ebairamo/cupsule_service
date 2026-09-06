package service

import (
	"capsule_service/internal/domain"
	"capsule_service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type CapsuleService struct {
	repo repository.CapsuleRepository
}

func NewCapsuleService(repo repository.CapsuleRepository) *CapsuleService {
	return &CapsuleService{
		repo: repo,
	}
}

func (s *CapsuleService) CreateCapsule(ctx context.Context, ownerID string, title string, message string, visibility domain.Visibility, openAt time.Time) (*domain.Cupsule, error) {
	if ownerID == "" || title == "" || message == "" || visibility == "" {
		return nil, domain.ErrValueEmpty
	}
	if openAt.Before(time.Now()) {
		return nil, domain.ErrInvalidOpenAt
	}

	cupsule := &domain.Cupsule{
		ID:         uuid.NewString(),
		OwnerID:    ownerID,
		Title:      title,
		Message:    message,
		Status:     domain.StatusDraft,
		Visibility: visibility,
		OpenAt:     openAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := s.repo.Create(ctx, cupsule); err != nil {
		return nil, err
	}
	return cupsule, nil
}

func (s *CapsuleService) GetByID(ctx context.Context, id string) (*domain.Cupsule, error) {
	if id == "" {
		return nil, domain.ErrValueEmpty
	}
	cupsule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return cupsule, nil
}
func (s *CapsuleService) List(ctx context.Context, ownerID string) ([]*domain.Cupsule, error) {
	if ownerID == "" {
		return nil, domain.ErrValueEmpty
	}
	cupsules, err := s.repo.List(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return cupsules, nil
}

func (s *CapsuleService) Update(ctx context.Context, cupsule domain.Cupsule) error {
	err := s.repo.Update(ctx, &cupsule)
	if err != nil {
		return err
	}
	return nil
}
func (s *CapsuleService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrValueEmpty
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
