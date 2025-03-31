package config

import (
	"os"
)

type Config struct {
	Database Database
	JWT      JWT
	Strava   Strava
	Summit   Summit
}

func NewConfig() *Config {
	return &Config{
		Database: Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBName:   os.Getenv("DATABASE_DBNAME"),
			SSLMode:  os.Getenv("DATABASE_SSLMODE"),
		},
		JWT: JWT{
			Secret: os.Getenv("JWT_SECRET"),
		},
		Strava: Strava{
			DistanceCacheTTL: os.Getenv("DISTANCE_CACHE_TTL"),
			ClientID:         os.Getenv("STRAVA_CLIENT_ID"),
			ClientSecret:     os.Getenv("STRAVA_CLIENT_SECRET"),
		},
		Summit: Summit{
			SummitThresholdMeters: os.Getenv("SUMMIT_THRESHOLD_METERS"),
		},
	}
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWT struct {
	Secret string
}

type Strava struct {
	DistanceCacheTTL string
	ClientID         string
	ClientSecret     string
}

type Summit struct {
	SummitThresholdMeters string // = "0.0007"
}
