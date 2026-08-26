package models

import "time"

// PageHit represents a page visit tracking record
type PageHit struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	UserID    *string   `json:"userId,omitempty"`
	UserEmail *string   `json:"userEmail,omitempty"`
	IPAddress *string   `json:"ipAddress,omitempty"`
	UserAgent *string   `json:"userAgent,omitempty"`
	Referrer  *string   `json:"referrer,omitempty"`
	Country   *string   `json:"country,omitempty"`
	Region    *string   `json:"region,omitempty"`
	City      *string   `json:"city,omitempty"`
	Latitude  *float64  `json:"latitude,omitempty"`
	Longitude *float64  `json:"longitude,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// ===== Request DTOs =====

// RecordHitRequest is the request body for recording a page hit
type RecordHitRequest struct {
	Path  string `json:"path" validate:"required"`
	Email string `json:"email,omitempty"`
}

// ===== Response DTOs =====

// PageHitStatsResponse is the aggregated hit statistics
type PageHitStatsResponse struct {
	Path       string `json:"path"`
	HitCount   int64  `json:"hitCount"`
	UniqueHits int64  `json:"uniqueHits"`
}

// PageHitSummaryResponse is the overall summary
type PageHitSummaryResponse struct {
	TotalHits int64                  `json:"totalHits"`
	TodayHits int64                  `json:"todayHits"`
	PageStats []PageHitStatsResponse `json:"pageStats"`
}

// PageHitDetailResponse is a single hit with geo info
type PageHitDetailResponse struct {
	ID        string   `json:"id"`
	Path      string   `json:"path"`
	UserEmail string   `json:"userEmail"`
	IPAddress *string  `json:"ipAddress,omitempty"`
	Country   *string  `json:"country,omitempty"`
	Region    *string  `json:"region,omitempty"`
	City      *string  `json:"city,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	CreatedAt string   `json:"createdAt"`
}

// HitsByCountryResponse is hit count grouped by country
type HitsByCountryResponse struct {
	Country  string `json:"country"`
	HitCount int64  `json:"hitCount"`
}

// HitsByRegionResponse is hit count grouped by country and region
type HitsByRegionResponse struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	HitCount int64  `json:"hitCount"`
}

// HitsHourlyResponse is hit count grouped by hour
type HitsHourlyResponse struct {
	Hour     string `json:"hour"`
	HitCount int64  `json:"hitCount"`
}

// HitsDailyResponse is hit count grouped by day
type HitsDailyResponse struct {
	Date     string `json:"date"`
	HitCount int64  `json:"hitCount"`
}

// ToDetailResponse converts a PageHit entity to a response DTO
func (h *PageHit) ToDetailResponse() *PageHitDetailResponse {
	email := "anonymous"
	if h.UserEmail != nil && *h.UserEmail != "" {
		email = *h.UserEmail
	}
	return &PageHitDetailResponse{
		ID:        h.ID,
		Path:      h.Path,
		UserEmail: email,
		IPAddress: h.IPAddress,
		Country:   h.Country,
		Region:    h.Region,
		City:      h.City,
		Latitude:  h.Latitude,
		Longitude: h.Longitude,
		CreatedAt: h.CreatedAt.Format(time.RFC3339),
	}
}
