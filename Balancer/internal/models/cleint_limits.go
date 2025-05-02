package models

import "context"

type ClientIdentifier string

const (
	IPPrefix     = "ip:"
	APIKeyPrefix = "api_key:"
)

type ClientLimit struct {
	ClientID string
	Capacity int
	Rate     float64
}

type Storage interface {
	GetClientLimit(ctx context.Context, clientID string) (*ClientLimit, error)
	SetClientLimit(ctx context.Context, clientID string, capacity int, rate float64) error
	GetAllClientLimits(ctx context.Context) ([]ClientLimit, error)
}
