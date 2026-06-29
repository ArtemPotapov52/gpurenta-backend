package db

import (
	"context"
	"fmt"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/types"
)

func (s *Store) FindUserByGoogleID(ctx context.Context, googleID string) (*types.User, error) {
	u := &types.User{}
	err := s.Pool.QueryRow(ctx,
		`SELECT id, google_id, email, name, picture, created_at FROM users WHERE google_id = $1`, googleID,
	).Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.Picture, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}

func (s *Store) CreateUser(ctx context.Context, googleID, email, name, picture string) (*types.User, error) {
	u := &types.User{}
	err := s.Pool.QueryRow(ctx,
		`INSERT INTO users (google_id, email, name, picture) VALUES ($1, $2, $3, $4)
		 RETURNING id, google_id, email, name, picture, created_at`,
		googleID, email, name, picture,
	).Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.Picture, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (s *Store) FindOrCreateUser(ctx context.Context, googleID, email, name, picture string) (*types.User, error) {
	u, err := s.FindUserByGoogleID(ctx, googleID)
	if err == nil {
		return u, nil
	}
	return s.CreateUser(ctx, googleID, email, name, picture)
}
