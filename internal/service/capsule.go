package service

import (
	"capsule_service/internal/repository"
	"context"
	"time"
)

type CapsuleService struct {
	repo repository.CapsuleRepository
}

func NewCapsuleService(repo repository.CapsuleRepository) *CapsuleService {
	return &CapsuleService{
		repo: repo,
	}
}

func (s *CapsuleService) CreateCapsule(ctx context.Context, ownerID string, title string, message string, visibility domain.Visibility, openAt time.Time) (*domain.Capsule, error) {

}
