package subscription

import (
	"context"
	"log/slog"

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

type Service struct {
	storage Storage
	logger  *slog.Logger
}

func NewService(s Storage, logger *slog.Logger) *Service {
	return &Service{storage: s, logger: logger}
}

func (s *Service) Create(ctx context.Context, sub *domain.Subscription) error {
	s.logger.Debug("service: create subscription", "service_name", sub.ServiceName, "user_id", sub.UserID.String())
	err := s.storage.Create(ctx, sub)
	if err != nil {
		s.logger.Error("service: failed to create subscription", "error", err)
		return err
	}
	s.logger.Info("service: subscription created", "subscription_id", sub.ID.String())
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	s.logger.Debug("service: get all subscriptions")
	subs, err := s.storage.GetAll(ctx)
	if err != nil {
		s.logger.Error("service: failed to get all subscriptions", "error", err)
		return nil, err
	}
	s.logger.Info("service: retrieved subscriptions", "count", len(subs))
	return subs, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	s.logger.Debug("service: get subscription by ID", "subscription_id", id.String())
	sub, err := s.storage.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("service: failed to get subscription by ID", "subscription_id", id.String(), "error", err)
		return nil, err
	}
	s.logger.Info("service: subscription retrieved", "subscription_id", id.String())
	return sub, nil
}

func (s *Service) Update(ctx context.Context, sub *domain.Subscription) error {
	s.logger.Debug("service: update subscription", "subscription_id", sub.ID.String())
	err := s.storage.Update(ctx, sub)
	if err != nil {
		s.logger.Error("service: failed to update subscription", "subscription_id", sub.ID.String(), "error", err)
		return err
	}
	s.logger.Info("service: subscription updated", "subscription_id", sub.ID.String())
	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("service: delete subscription", "subscription_id", id.String())
	err := s.storage.Delete(ctx, id)
	if err != nil {
		s.logger.Error("service: failed to delete subscription", "subscription_id", id.String(), "error", err)
		return err
	}
	s.logger.Info("service: subscription deleted", "subscription_id", id.String())
	return nil
}

func (s *Service) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to domain.YearMonth) (int64, error) {
	s.logger.Debug("service: calculate total cost",
		"user_id", userID,
		"service_name", serviceName,
		"from", from.Time,
		"to", to.Time,
	)
	total, err := s.storage.TotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		s.logger.Error("service: failed to calculate total cost", "error", err)
		return 0, err
	}
	s.logger.Info("service: total cost calculated", "total", total)
	return total, nil
}
