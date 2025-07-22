package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRequestDTO struct {
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"499"`
	UserID      string  `json:"user_id" example:"e5c7c66b-4a3e-4728-84d9-b6c6b46ef1a6"`
	StartDate   string  `json:"start_date" example:"07-2024"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2024"`
}

type SubscriptionResponseDTO struct {
	ID          string  `json:"id" example:"696c530f-b6c5-467f-ab70-45916e72daa7"`
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"499"`
	UserID      string  `json:"user_id" example:"e5c7c66b-4a3e-4728-84d9-b6c6b46ef1a6"`
	StartDate   string  `json:"start_date" example:"07-2024"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2024"`
}

func (dto *SubscriptionRequestDTO) Validate() error {
	if strings.TrimSpace(dto.ServiceName) == "" {
		return errors.New("service_name is required")
	}
	if dto.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if _, err := uuid.Parse(dto.UserID); err != nil {
		return errors.New("user_id is invalid UUID")
	}
	if _, err := time.Parse("01-2006", dto.StartDate); err != nil {
		return errors.New("start_date has invalid format, expected MM-YYYY")
	}
	if dto.EndDate != nil {
		endTime, err := time.Parse("01-2006", *dto.EndDate)
		if err != nil {
			return errors.New("end_date has invalid format, expected MM-YYYY")
		}
		startTime, _ := time.Parse("01-2006", dto.StartDate)
		if endTime.Before(startTime) {
			return errors.New("end_date cannot be before start_date")
		}
	}
	return nil
}