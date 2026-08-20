package repository

import (
	"capsule_service/internal/domain"
	"context"
)

type CapsuleRepository interface {
	Create(ctx context.Context, capsule *domain.Cupsule) error
	GetByID(ctx context.Context, id string) (*domain.Cupsule, error)
	List(ctx context.Context, ownerID string) ([]*domain.Cupsule, error)
	Update(ctx context.Context, capsule *domain.Cupsule) error
	Delete(ctx context.Context, id string) error
}
