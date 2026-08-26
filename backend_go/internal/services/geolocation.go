package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"menunderfire/internal/logger"
	"menunderfire/internal/repositories"
)

// GeoResult represents the response from ip-api.com
type GeoResult struct {
	Status  string  `json:"status"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// LookupAndUpdateGeo performs an IP geolocation lookup and updates the hit record.
// This runs asynchronously (fire-and-forget) using ip-api.com free tier.
func LookupAndUpdateGeo(hitRepo repositories.PageHitRepository, hitID, ipAddress string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,city,lat,lon", ipAddress)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			logger.Debug().Err(err).Str("ip", ipAddress).Msg("Failed to create geo request")
			return
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logger.Debug().Err(err).Str("ip", ipAddress).Msg("Failed to lookup geo")
			return
		}
		defer func() { _ = resp.Body.Close() }()

		var geo GeoResult
		if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
			logger.Debug().Err(err).Str("ip", ipAddress).Msg("Failed to decode geo response")
			return
		}

		if geo.Status != "success" {
			logger.Debug().Str("ip", ipAddress).Str("status", geo.Status).Msg("Geo lookup failed")
			return
		}

		var country, city *string
		var lat, lon *float64
		if geo.Country != "" {
			country = &geo.Country
		}
		if geo.City != "" {
			city = &geo.City
		}
		if geo.Lat != 0 || geo.Lon != 0 {
			lat = &geo.Lat
			lon = &geo.Lon
		}

		if err := hitRepo.UpdateGeo(context.Background(), hitID, country, city, lat, lon); err != nil {
			logger.Debug().Err(err).Str("hit_id", hitID).Msg("Failed to update geo data")
			return
		}

		logger.Debug().
			Str("hit_id", hitID).
			Str("country", geo.Country).
			Str("city", geo.City).
			Msg("Geo data updated for hit")
	}()
}
