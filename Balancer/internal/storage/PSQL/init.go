package PSQL

import (
	"context"
)

func (s *PSQLStorage) InitSchema(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS client_limits (
            client_id TEXT PRIMARY KEY,
            capacity INTEGER NOT NULL,
            rate FLOAT NOT NULL,
            updated_at TIMESTAMP DEFAULT NOW()
        )`)
	return err
}
