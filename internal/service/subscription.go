package service

import (
	"context"

	"subscription-service/internal/domain"

	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetAll(ctx context.Context) ([]*domain.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to domain.YearMonth) (int64, error)
}

type SubscriptionService struct {
	storage Storage
}

func NewSubscriptionService(s Storage) *SubscriptionService {
	return &SubscriptionService{storage: s}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *domain.Subscription) error {
	return s.storage.Create(ctx, sub)
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	return s.storage.GetAll(ctx)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *domain.Subscription) error {
	return s.storage.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.storage.Delete(ctx, id)
}

func (s *SubscriptionService) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to domain.YearMonth) (int64, error) {
	return s.storage.TotalCost(ctx, userID, serviceName, from, to)
}
