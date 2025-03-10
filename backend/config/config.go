package config

import (
	"os"
)

type Config struct {
	Strava   Strava
	Database Database
	Summit   Summit
}

func NewConfig() *Config {
	return &Config{
		Strava: Strava{
			DistanceCacheTTL: os.Getenv("DISTANCE_CACHE_TTL"),
			ClientID:         os.Getenv("STRAVA_CLIENT_ID"),
			ClientSecret:     os.Getenv("STRAVA_CLIENT_SECRET"),
		},
		Database: Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBName:   os.Getenv("DATABASE_DBNAME"),
			SSLMode:  os.Getenv("DATABASE_SSLMODE"),
		},
		Summit: Summit{
			SummitThresholdMeters: os.Getenv("SUMMIT_THRESHOLD_METERS"),
		},
	}
}

type Strava struct {
	DistanceCacheTTL string
	ClientID         string
	ClientSecret     string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Summit struct {
	SummitThresholdMeters string // = "0.0007"
}
