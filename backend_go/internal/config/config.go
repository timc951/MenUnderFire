package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	CORSOrigin          string
	HitToken            string
	AgreementHMACSecret string
	AuthDomain          string
	AuthAudience        string
	AuthIssuer          string
	TrustProxyHeaders   bool
}

func Load() *Config {
	return &Config{
		ServerPort:          getEnv("SERVER_PORT", "7001"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              requireEnv("DB_USER"),
		DBPassword:          requireEnv("DB_PASSWORD"),
		DBName:              getEnv("DB_NAME", "menunderfire"),
		DBSSLMode:           getEnv("DB_SSL_MODE", "require"),
		CORSOrigin:          requireEnv("CORS_ALLOWED_ORIGIN"),
		HitToken:            getEnv("HIT_TOKEN", ""),
		AgreementHMACSecret: requireEnv("AGREEMENT_HMAC_SECRET"),
		AuthDomain:          getEnv("AUTH_DOMAIN", ""),
		AuthAudience:        getEnv("AUTH_AUDIENCE", ""),
		AuthIssuer:          getEnv("AUTH_ISSUER", ""),
		TrustProxyHeaders:   getEnv("TRUST_PROXY_HEADERS", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func requireEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}
