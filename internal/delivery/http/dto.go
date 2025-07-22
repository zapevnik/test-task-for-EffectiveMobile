package http

import (
	"subscription-service/internal/domain"
	"github.com/google/uuid"
)

type CreateSubscriptionDTO struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   domain.YearMonth  `json:"start_date"`
	EndDate     *domain.YearMonth `json:"end_date,omitempty"`
}