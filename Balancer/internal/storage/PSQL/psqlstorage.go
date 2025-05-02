package PSQL

import (
	"cloud/Balancer/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PSQLStorage struct {
	db *sql.DB
}

func NewStorage(connStr string) (*PSQLStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PSQLStorage{db: db}, nil
}

func (s *PSQLStorage) GetClientLimit(ctx context.Context, clientID string) (*models.ClientLimit, error) {
	var limit models.ClientLimit
	err := s.db.QueryRowContext(ctx,
		"SELECT client_id, capacity, rate FROM client_limits WHERE client_id = $1", // запрос id, cap, rate
		clientID).Scan(&limit.ClientID, &limit.Capacity, &limit.Rate)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (s *PSQLStorage) SetClientLimit(ctx context.Context, clientID string, capacity int, rate float64) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO client_limits (client_id, capacity, rate) 
		VALUES ($1, $2, $3)
		ON CONFLICT (client_id) 
		DO UPDATE SET capacity = $2, rate = $3, updated_at = NOW()`,
		clientID, capacity, rate)
	return err
}

func (s *PSQLStorage) GetAllClientLimits(ctx context.Context) ([]models.ClientLimit, error) {
	const query = `SELECT client_id, capacity, rate FROM client_limits`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var limits []models.ClientLimit
	for rows.Next() {
		var l models.ClientLimit
		if err := rows.Scan(&l.ClientID, &l.Capacity, &l.Rate); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		limits = append(limits, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return limits, nil
}
