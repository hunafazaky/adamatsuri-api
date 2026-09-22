package config

import (
	"fmt"
	"os"
)

// Config holds every environment-derived setting the app needs, loaded once
// at startup. Passing this struct down (instead of calling os.Getenv all
// over the codebase) makes it obvious what each package actually depends on,
// and makes it possible to construct a Config by hand in tests.
type Config struct {
	Port string

	DBURI        string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresSSL  string

	JWTSecret          string
	ImageKitPrivateKey string

	// ClientOrigin is the frontend's origin(s), comma-separated, allowed
	// to call this API from a browser (CORS). Defaults to the Vite dev
	// server's default port so local dev works with zero setup.
	ClientOrigin string

	// PublicHost is the externally-reachable host:port for this deployment
	// (e.g. "your-app.onrender.com"). Optional — empty locally, where the
	// spec's default @host (localhost:8080) from the swag annotation
	// already suffices. When set, main.go uses it to override
	// docs.SwaggerInfo.Host so "try it out" in Swagger UI / Scalar hits
	// the right domain instead of localhost.
	PublicHost string
}

// Load reads environment variables into a Config and validates that
// everything required is present. Call this once in main() — if it fails,
// fail fast at startup rather than hundreds of requests in.
func Load() (*Config, error) {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),

		DBURI:        os.Getenv("DB_URI"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:   os.Getenv("POSTGRES_DB_NAME"),
		PostgresSSL:  os.Getenv("POSTGRES_SSLMODE"),

		JWTSecret:          os.Getenv("JWT_SECRET"),
		ImageKitPrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),

		ClientOrigin: getEnv("CLIENT_ORIGIN", "http://localhost:5173"),

		// PUBLIC_HOST is the explicit override. If not set, fall back to
		// RENDER_EXTERNAL_HOSTNAME, which Render auto-populates for every
		// web service — meaning on Render, this needs NO manual setup at
		// all. PUBLIC_HOST still exists for non-Render deployments, or if
		// you later attach a custom domain and want to override it.
		PublicHost: firstNonEmpty(os.Getenv("PUBLIC_HOST"), os.Getenv("RENDER_EXTERNAL_HOSTNAME")),
	}

	if cfg.DBURI == "" {
		missing := []string{}
		if cfg.PostgresHost == "" {
			missing = append(missing, "POSTGRES_HOST")
		}
		if cfg.PostgresPort == "" {
			missing = append(missing, "POSTGRES_PORT")
		}
		if cfg.PostgresUser == "" {
			missing = append(missing, "POSTGRES_USER")
		}
		if cfg.PostgresPass == "" {
			missing = append(missing, "POSTGRES_PASSWORD")
		}
		if cfg.PostgresDB == "" {
			missing = append(missing, "POSTGRES_DB_NAME")
		}
		if cfg.PostgresSSL == "" {
			missing = append(missing, "POSTGRES_SSLMODE")
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("missing required environment variables (or set DB_URI instead): %v", missing)
		}
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// DSN builds the Postgres connection string, preferring DB_URI when set.
func (c *Config) DSN() string {
	if c.DBURI != "" {
		return c.DBURI
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.PostgresHost, c.PostgresUser, c.PostgresPass, c.PostgresDB, c.PostgresPort, c.PostgresSSL,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// firstNonEmpty returns the first non-empty string, or "" if all are empty.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
