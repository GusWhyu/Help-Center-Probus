package config

import "testing"

func TestLoadConfigDatabaseDrivers(t *testing.T) {
	for key, value := range map[string]string{
		"DB_DRIVER":        "  MySQL ",
		"DB_HOST":          "",
		"DB_PORT":          "",
		"DB_USER":          "",
		"DB_PASSWORD":      "",
		"DB_USER_NAME":     "secondary",
		"DB_USER_DRIVER":   "postgres",
		"DB_USER_HOST":     "",
		"DB_USER_PORT":     "",
		"DB_USER_USER":     "",
		"DB_USER_PASSWORD": "",
		"DB_USER_SSLMODE":  "",
	} {
		t.Setenv(key, value)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	primary := cfg.Databases["default"]
	if primary.Driver != "mysql" || primary.Port != "3306" || primary.User != "root" {
		t.Errorf("unexpected MySQL defaults: driver=%q port=%q user=%q", primary.Driver, primary.Port, primary.User)
	}

	secondary := cfg.Databases["user"]
	if !secondary.Optional || secondary.Driver != "postgres" || secondary.Port != "5432" || secondary.User != "postgres" {
		t.Errorf("unexpected PostgreSQL secondary configuration: optional=%t driver=%q port=%q user=%q",
			secondary.Optional, secondary.Driver, secondary.Port, secondary.User)
	}
}

func TestLoadConfigOmitsOptionalUserDatabaseWhenNameIsEmpty(t *testing.T) {
	t.Setenv("DB_USER_NAME", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if _, exists := cfg.Databases["user"]; exists {
		t.Fatal("optional user database should not be configured when DB_USER_NAME is empty")
	}
}

func TestLoadConfigRequiresStrongJWTSecretInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "short-secret")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("LoadConfig() error = nil, want production JWT secret validation error")
	}

	t.Setenv("JWT_SECRET", "production-secret-with-at-least-32-characters")
	if _, err := LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() with strong secret error = %v", err)
	}
}
