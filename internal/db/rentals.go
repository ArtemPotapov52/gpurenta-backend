package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/types"
)

func (s *Store) CreateRental(ctx context.Context, agentID, renterID, image, frpURL string, hours int) (*types.Rental, error) {
	r := &types.Rental{}
	now := time.Now()
	endsAt := now.Add(time.Duration(hours) * time.Hour)
	err := s.Pool.QueryRow(ctx,
		`INSERT INTO rentals (agent_id, renter_id, image, frp_url, ends_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, agent_id, renter_id, image, frp_url, cost_cents, status, started_at, ends_at`,
		agentID, renterID, image, frpURL, endsAt,
	).Scan(&r.ID, &r.AgentID, &r.RenterID, &r.Image, &r.FrpURL, &r.CostCents, &r.Status, &r.StartedAt, &r.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("create rental: %w", err)
	}
	return r, nil
}

func (s *Store) GetRentalByID(ctx context.Context, id string) (*types.Rental, error) {
	r := &types.Rental{}
	err := s.Pool.QueryRow(ctx,
		`SELECT id, agent_id, renter_id, image, frp_url, cost_cents, status, started_at, ends_at
		 FROM rentals WHERE id = $1`, id,
	).Scan(&r.ID, &r.AgentID, &r.RenterID, &r.Image, &r.FrpURL, &r.CostCents, &r.Status, &r.StartedAt, &r.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("get rental: %w", err)
	}
	return r, nil
}

func (s *Store) StopRental(ctx context.Context, id string) (*types.Rental, error) {
	r := &types.Rental{}
	now := time.Now()
	err := s.Pool.QueryRow(ctx,
		`UPDATE rentals SET status = 'completed', ends_at = $2,
		 cost_cents = EXTRACT(EPOCH FROM $2 - started_at)::int / 36
		 WHERE id = $1 AND status = 'active'
		 RETURNING id, agent_id, renter_id, image, frp_url, cost_cents, status, started_at, ends_at`,
		id, now,
	).Scan(&r.ID, &r.AgentID, &r.RenterID, &r.Image, &r.FrpURL, &r.CostCents, &r.Status, &r.StartedAt, &r.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("stop rental: %w", err)
	}
	return r, nil
}

func (s *Store) GetActiveRentalByAgentID(ctx context.Context, agentID string) (*types.Rental, error) {
	r := &types.Rental{}
	err := s.Pool.QueryRow(ctx,
		`SELECT id, agent_id, renter_id, image, frp_url, cost_cents, status, started_at, ends_at
		 FROM rentals WHERE agent_id = $1 AND status = 'active' LIMIT 1`, agentID,
	).Scan(&r.ID, &r.AgentID, &r.RenterID, &r.Image, &r.FrpURL, &r.CostCents, &r.Status, &r.StartedAt, &r.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("get active rental: %w", err)
	}
	return r, nil
}
