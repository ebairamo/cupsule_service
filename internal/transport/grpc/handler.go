package grpc

import (
	"context"
	"errors"

	"capsule_service/internal/domain"
	service "capsule_service/internal/usecase"

	capsulev1 "github.com/ebairamo/digital-capsule-contracts-go/gen/go/capsule/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	capsulev1.UnimplementedCapsuleServiceServer
	capsuleService *service.CapsuleService
}

func NewHandler(capsuleService *service.CapsuleService) *Handler {
	return &Handler{
		capsuleService: capsuleService,
	}
}

func (h *Handler) GetCapsule(
	ctx context.Context,
	req *capsulev1.GetCapsuleRequest,
) (*capsulev1.Capsule, error) {
	// 1. Валидация входных данных
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 2. Вызов бизнес-логики из слоя service
	capsule, err := h.capsuleService.GetByID(ctx, req.GetId())
	if err != nil {
		// Ошибку "не найдено" маппим в codes.NotFound
		if errors.Is(err, domain.ErrCapsuleNotFound) { // укажи свою ошибку из domain/service
			return nil, status.Error(codes.NotFound, "capsule not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	// 3. Маппинг сущности domain.Capsule в Protobuf DTO
	return &capsulev1.Capsule{
		Id:         capsule.ID,
		OwnerId:    capsule.OwnerID,
		Title:      capsule.Title,
		Message:    capsule.Message,
		Status:     capsulev1.Status(capsule.Status),
		Visibility: capsulev1.Visibility(capsule.Visibility),
		OpenAt:     timestamppb.New(capsule.OpenAt),
		UpdatedAt:  timestamppb.New(capsule.UpdatedAt),
	}, nil

}
