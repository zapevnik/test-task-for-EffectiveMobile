package dtoConv

import (
	"errors"
	"subscription-service/internal/delivery/dto"
	"subscription-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

func  RequestDtoToDomain(res dto.SubscriptionRequestDTO) (*domain.Subscription, error) {
	userID, err := uuid.Parse(res.UserID)
	if err != nil {
		return nil, errors.New("invalid user_id")
	}

	startTime, err := time.Parse("01-2006", res.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format")
	}

	var endYearMonth *domain.YearMonth
	if res.EndDate != nil {
		endTime, err := time.Parse("01-2006", *res.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format")
		}
		ym := domain.YearMonth{Time: endTime}
		endYearMonth = &ym
	}

	sub := &domain.Subscription{
		ID:          uuid.Nil,
		ServiceName: res.ServiceName,
		Price:       res.Price,
		UserID:      userID,
		StartDate:   domain.YearMonth{Time: startTime},
		EndDate:     endYearMonth,
	}

	return sub, nil
}

func DomainToResponseDTO(sub *domain.Subscription) dto.SubscriptionResponseDTO {
    var endDate *string
    if sub.EndDate != nil {
        s := sub.EndDate.Format("01-2006")
        endDate = &s
    }
    return dto.SubscriptionResponseDTO{
        ID:          sub.ID.String(),
        ServiceName: sub.ServiceName,
        Price:       sub.Price,
        UserID:      sub.UserID.String(),
        StartDate:   sub.StartDate.Format("01-2006"),
        EndDate:     endDate,
    }
}
