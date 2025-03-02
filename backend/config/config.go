package config

import (
	"os"
	"strconv"
)

type Config struct {
	Strava Strava
}

func NewConfig() *Config {
	return &Config{
		Strava: Strava {
			DistanceCacheTTL: os.Getenv("DISTANCE_CACHE_TTL"),
			ClientID: os.Getenv("STRAVA_CLIENT_ID"),
			ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		},
	}
}

type Strava struct {
	DistanceCacheTTL string
	ClientID string
	ClientSecret string
}
