package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"subscription-service/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionStorage struct {
	db *sql.DB
}

func NewSubscriptionStorage(db *sql.DB) *SubscriptionStorage {
	return &SubscriptionStorage{db: db}
}

func (s *SubscriptionStorage) Create(ctx context.Context, sub *domain.Subscription) error {
	sub.ID = uuid.New()

	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	return err
}

func (s *SubscriptionStorage) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		sub := new(domain.Subscription)
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SubscriptionStorage) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	sub := new(domain.Subscription)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionStorage) Update(ctx context.Context, sub *domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE id = $5
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	return err
}

func (s *SubscriptionStorage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *SubscriptionStorage) TotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	from, to time.Time,
) (int64, error) {
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
	return total, err
}
