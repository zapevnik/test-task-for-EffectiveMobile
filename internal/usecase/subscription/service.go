package subscription

import (
	"context"
	"time"

	"subscription-service/internal/domain"
	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetAll(ctx context.Context) ([]*domain.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int64, error)
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) Create(ctx context.Context, sub *domain.Subscription) error {
	return s.storage.Create(ctx, sub)
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	return s.storage.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, sub *domain.Subscription) error {
	return s.storage.Update(ctx, sub)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.storage.Delete(ctx, id)
}

func (s *Service) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int64, error) {
	return s.storage.TotalCost(ctx, userID, serviceName, from, to)
}
