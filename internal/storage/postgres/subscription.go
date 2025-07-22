package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"subscription-service/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionStorage struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSubscriptionStorage(db *sql.DB, logger *slog.Logger) *SubscriptionStorage {
	return &SubscriptionStorage{db: db, logger: logger}
}

func (s *SubscriptionStorage) Create(ctx context.Context, sub *domain.Subscription) error {
	sub.ID = uuid.New()
	s.logger.Info("Create subscription started", "id", sub.ID.String(), "service_name", sub.ServiceName, "user_id", sub.UserID.String())

	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	var endDate *time.Time
	if sub.EndDate != nil {
		t := sub.EndDate.Time
		endDate = &t
	}
	_, err := s.db.ExecContext(
		ctx,
		query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate.Time,
		endDate,
	)
	if err != nil {
		s.logger.Error("Create subscription failed", "id", sub.ID.String(), "error", err)
		return err
	}

	s.logger.Info("Create subscription succeeded", "id", sub.ID.String())
	return nil
}

func (s *SubscriptionStorage) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	s.logger.Info("GetAll subscriptions started")

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		s.logger.Error("GetAll subscriptions query failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		var (
			start time.Time
			end   *time.Time
		)

		sub := new(domain.Subscription)
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&start,
			&end,
		)
		if err != nil {
			s.logger.Error("GetAll subscriptions scan failed", "error", err)
			return nil, err
		}

		sub.StartDate.Time = start
		if end != nil {
			sub.EndDate = &domain.YearMonth{Time: *end}
		}

		subs = append(subs, sub)
	}

	s.logger.Info("GetAll subscriptions succeeded", "count", len(subs))
	return subs, nil
}


func (s *SubscriptionStorage) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	s.logger.Info("GetByID subscription started", "id", id.String())

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	var (
		start time.Time
		end   *time.Time
	)

	sub := new(domain.Subscription)

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&start,
		&end,
	)
	if err != nil {
		s.logger.Error("GetByID subscription failed", "id", id.String(), "error", err)
		return nil, err
	}

	sub.StartDate.Time = start

	if end != nil {
		sub.EndDate = &domain.YearMonth{Time: *end}
	}

	s.logger.Info("GetByID subscription succeeded", "id", id.String())
	return sub, nil
}


func (s *SubscriptionStorage) Update(ctx context.Context, sub *domain.Subscription) error {
	s.logger.Info("Update subscription started", "id", sub.ID.String())

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE id = $5
	`

	var endDate *time.Time
	if sub.EndDate != nil {
		t := sub.EndDate.Time
		endDate = &t
	}

	_, err := s.db.ExecContext(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate.Time,
		endDate,
		sub.ID,
	)
	if err != nil {
		s.logger.Error("Update subscription failed", "id", sub.ID.String(), "error", err)
		return err
	}

	s.logger.Info("Update subscription succeeded", "id", sub.ID.String())
	return nil
}

func (s *SubscriptionStorage) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Delete subscription started", "id", id.String())

	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		s.logger.Error("Delete subscription failed", "id", id.String(), "error", err)
		return err
	}

	s.logger.Info("Delete subscription succeeded", "id", id.String())
	return nil
}

func (s *SubscriptionStorage) TotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from, to time.Time,
) (int64, error) {
	s.logger.Info("TotalCost calculation started", "from", from.Format("01-2006"), "to", to.Format("01-2006"))
	if userID != nil {
		s.logger.Info("TotalCost filter by userID", "userID", userID.String())
	}
	if serviceName != nil {
		s.logger.Info("TotalCost filter by serviceName", "serviceName", *serviceName)
	}

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date >= $1 AND start_date <= $2
	`
	args := []interface{}{from, to}
	argIdx := 3

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, *userID)
		argIdx++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, *serviceName)
	}

	var total int64
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		s.logger.Error("TotalCost calculation failed", "error", err)
		return 0, err
	}

	s.logger.Info("TotalCost calculation succeeded", "total", total)
	return total, nil
}
