package config

import "os"

const (
	// LogLevel const
	LogLevel                = "info"
	secretGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
	goEnvironment           = "GO_ENVIRONMENT"
	production              = "production"
)

var githubAccessToken = os.Getenv(secretGithubAccessToken)

// GetGithubAccessToken func
func GetGithubAccessToken() string {
	return githubAccessToken
}

// IsProduction func
func IsProduction() bool {
	return os.Getenv(goEnvironment) == production
}
