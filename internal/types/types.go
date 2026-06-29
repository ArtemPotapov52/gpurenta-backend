package types

import "time"

type User struct {
	ID        string    `json:"id"`
	GoogleID  string    `json:"-"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Agent struct {
	ID              string     `json:"id"`
	OwnerID         string     `json:"owner_id"`
	GPUModel        string     `json:"gpu_model"`
	VRAMGB          int        `json:"vram_gb"`
	OS              string     `json:"os"`
	FRPURL          string     `json:"frp_url,omitempty"`
	Status          string     `json:"status"`
	SupportedImages []string   `json:"supported_images"`
	PricePerHour    int        `json:"price_per_hour"`
	Secret          string     `json:"-"`
	LastHeartbeat   *time.Time `json:"last_heartbeat,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Rental struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	RenterID  string    `json:"renter_id"`
	Image     string    `json:"image"`
	FrpURL    string    `json:"frp_url,omitempty"`
	CostCents int       `json:"cost_cents,omitempty"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	EndsAt    time.Time `json:"ends_at,omitempty"`
}

type WorkloadImage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Port        int    `json:"port"`
}

var SupportedImages = []WorkloadImage{
	{Name: "ollama", Description: "LLM inference (llama, qwen, mistral...)", Port: 11434},
	{Name: "comfyui", Description: "Stable Diffusion / Flux with nodes", Port: 8188},
	{Name: "whisper", Description: "Speech recognition", Port: 9000},
	{Name: "sd", Description: "Stable Diffusion image generation", Port: 7860},
}
