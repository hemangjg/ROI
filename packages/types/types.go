package types

// Environment names used across the platform.
const (
	EnvLocal       = "local"
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// Pagination is a foundation placeholder for list endpoints (Phase 2+).
type Pagination struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}