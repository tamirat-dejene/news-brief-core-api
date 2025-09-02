package config

import (
	"os"
	"strconv"
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
)

// Config holds application configuration values.
type Config struct {
	SendActivationEmail          bool
	AppBaseURL                   string
	FrontendBaseURL              string
	FrontendMobileBaseURL        string
	RefreshTokenExpiry           time.Duration
	PasswordResetTokenExpiry     time.Duration
	EmailVerificationTokenExpiry time.Duration
}

// NewConfig creates a new Config instance, loading values from environment variables.
func NewConfig() contract.IConfigProvider {
	return &Config{
		SendActivationEmail:          getEnvAsBool("SEND_ACTIVATION_EMAIL", false),
		AppBaseURL:                   getEnv("APP_BASE_URL", "http://localhost:8080"),
		FrontendBaseURL:              getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
		FrontendMobileBaseURL:        getEnv("FRONTEND_MOBILE_BASE_URL", ""),
		RefreshTokenExpiry:           time.Hour * time.Duration(getEnvAsInt("REFRESH_TOKEN_EXPIRY_HOURS", 168)), // 7 days
		PasswordResetTokenExpiry:     time.Minute * time.Duration(getEnvAsInt("PASSWORD_RESET_TOKEN_EXPIRY_MINUTES", 15)),
		EmailVerificationTokenExpiry: time.Minute * time.Duration(getEnvAsInt("EMAIL_VERIFICATION_TOKEN_EXPIRY_MINUTES", 60)),
	}
}

// GetSendActivationEmail returns whether to send an activation email.
func (c *Config) GetSendActivationEmail() bool {
	return c.SendActivationEmail
}

// GetAppBaseURL returns the base URL of the application.
func (c *Config) GetAppBaseURL() string {
	return c.AppBaseURL
}

// GetFrontendBaseURL returns the frontend application's base URL.
func (c *Config) GetFrontendBaseURL() string {
	return c.FrontendBaseURL
}

// GetFrontendMobileBaseURL returns the mobile app's deep link/scheme base URL.
func (c *Config) GetFrontendMobileBaseURL() string {
	return c.FrontendMobileBaseURL
}

// GetRefreshTokenExpiry returns the expiry duration for refresh tokens.
func (c *Config) GetRefreshTokenExpiry() time.Duration {
	return c.RefreshTokenExpiry
}

// GetPasswordResetTokenExpiry returns the expiry duration for password reset tokens.
func (c *Config) GetPasswordResetTokenExpiry() time.Duration {
	return c.PasswordResetTokenExpiry
}

// GetEmailVerificationTokenExpiry returns the expiry duration for email verification tokens.
func (c *Config) GetEmailVerificationTokenExpiry() time.Duration {
	return c.EmailVerificationTokenExpiry
}

// Helper function to get an environment variable or return a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper function to get an environment variable as an integer or return a default value.
func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

// Helper function to get an environment variable as a boolean or return a default value.
func getEnvAsBool(name string, fallback bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}

func (c *Config) GetAIServiceAPIKey() string {
	return getEnv("AI_SERVICE_API_KEY", "")
}
