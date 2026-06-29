package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtemPotapov52/gpurenta-backend/internal/types"
)

func (s *Store) CreateAgent(ctx context.Context, ownerID, gpuModel, os string, vramGB int, images []string, price int) (*types.Agent, error) {
	a := &types.Agent{}
	err := s.Pool.QueryRow(ctx,
		`INSERT INTO agents (owner_id, gpu_model, vram_gb, os, supported_images, price_per_hour)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, owner_id, gpu_model, vram_gb, os, COALESCE(frp_url, ''), status, supported_images, price_per_hour, secret, last_heartbeat, created_at`,
		ownerID, gpuModel, vramGB, os, images, price,
	).Scan(&a.ID, &a.OwnerID, &a.GPUModel, &a.VRAMGB, &a.OS, &a.FRPURL, &a.Status, &a.SupportedImages, &a.PricePerHour, &a.Secret, &a.LastHeartbeat, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}
	return a, nil
}

func (s *Store) GetAgentBySecret(ctx context.Context, agentID, secret string) (*types.Agent, error) {
	a := &types.Agent{}
	err := s.Pool.QueryRow(ctx,
		`SELECT id, owner_id, gpu_model, vram_gb, os, COALESCE(frp_url, ''), status, supported_images, price_per_hour, secret, last_heartbeat, created_at
		 FROM agents WHERE id = $1 AND secret = $2`, agentID, secret,
	).Scan(&a.ID, &a.OwnerID, &a.GPUModel, &a.VRAMGB, &a.OS, &a.FRPURL, &a.Status, &a.SupportedImages, &a.PricePerHour, &a.Secret, &a.LastHeartbeat, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get agent by secret: %w", err)
	}
	return a, nil
}

func (s *Store) GetAgentByID(ctx context.Context, id string) (*types.Agent, error) {
	a := &types.Agent{}
	err := s.Pool.QueryRow(ctx,
		`SELECT id, owner_id, gpu_model, vram_gb, os, COALESCE(frp_url, ''), status, supported_images, price_per_hour, secret, last_heartbeat, created_at
		 FROM agents WHERE id = $1`, id,
	).Scan(&a.ID, &a.OwnerID, &a.GPUModel, &a.VRAMGB, &a.OS, &a.FRPURL, &a.Status, &a.SupportedImages, &a.PricePerHour, &a.Secret, &a.LastHeartbeat, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get agent: %w", err)
	}
	return a, nil
}

func (s *Store) Heartbeat(ctx context.Context, agentID, frpURL string) error {
	now := time.Now()
	_, err := s.Pool.Exec(ctx,
		`UPDATE agents SET status = 'online', last_heartbeat = $2, frp_url = COALESCE(NULLIF($3, ''), frp_url) WHERE id = $1`,
		agentID, now, frpURL,
	)
	return err
}

func (s *Store) ListOnlineGPUs(ctx context.Context, minVRAM int, imageFilter string) ([]types.Agent, error) {
	q := `SELECT id, owner_id, gpu_model, vram_gb, os, COALESCE(frp_url, ''), status, supported_images, price_per_hour, secret, last_heartbeat, created_at
		  FROM agents WHERE status = 'online'
		  AND (SELECT COUNT(*) FROM rentals WHERE agent_id = agents.id AND status = 'active') = 0`
	args := []interface{}{}
	argIdx := 1

	if minVRAM > 0 {
		argIdx++
		q += fmt.Sprintf(" AND vram_gb >= $%d", argIdx)
		args = append(args, minVRAM)
	}
	if imageFilter != "" {
		argIdx++
		q += fmt.Sprintf(" AND $%d = ANY(supported_images)", argIdx)
		args = append(args, imageFilter)
	}
	q += " ORDER BY price_per_hour ASC"

	rows, err := s.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list gpus: %w", err)
	}
	defer rows.Close()

	var agents []types.Agent
	for rows.Next() {
		var a types.Agent
		if err := rows.Scan(&a.ID, &a.OwnerID, &a.GPUModel, &a.VRAMGB, &a.OS, &a.FRPURL, &a.Status, &a.SupportedImages, &a.PricePerHour, &a.Secret, &a.LastHeartbeat, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan agent: %w", err)
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (s *Store) MarkAgentOffline(ctx context.Context, agentID string) error {
	_, err := s.Pool.Exec(ctx, `UPDATE agents SET status = 'offline' WHERE id = $1`, agentID)
	return err
}

func (s *Store) MarkStaleAgentsOffline(ctx context.Context, timeout time.Duration) error {
	_, err := s.Pool.Exec(ctx,
		`UPDATE agents SET status = 'offline' WHERE status = 'online' AND last_heartbeat < NOW() - $1::interval`,
		timeout.String(),
	)
	return err
}
