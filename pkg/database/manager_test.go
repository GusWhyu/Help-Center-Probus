package database

import (
	"testing"

	"github.com/rafialfauzan/layrd/pkg/config"
)

func TestDatabaseDialector(t *testing.T) {
	tests := []struct {
		name       string
		driver     string
		wantDriver string
	}{
		{name: "postgres", driver: "postgres", wantDriver: "postgres"},
		{name: "mysql", driver: "mysql", wantDriver: "mysql"},
		{name: "empty driver defaults to postgres", driver: "", wantDriver: "postgres"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialector, err := databaseDialector(config.DatabaseConfig{
				Name:     "test",
				Driver:   tt.driver,
				Host:     "127.0.0.1",
				Port:     "5432",
				User:     "user",
				Password: "pass",
				DBName:   "test",
				SSLMode:  "disable",
			}, "host=127.0.0.1")
			if err != nil {
				t.Fatalf("databaseDialector() error = %v", err)
			}
			if got := dialector.Name(); got != tt.wantDriver {
				t.Errorf("databaseDialector().Name() = %q, want %q", got, tt.wantDriver)
			}
		})
	}
}

func TestDatabaseDialectorRejectsUnsupportedDriver(t *testing.T) {
	_, err := databaseDialector(config.DatabaseConfig{
		Name:   "test",
		Driver: "sqlite",
	}, "")
	if err == nil {
		t.Fatal("databaseDialector() error = nil, want unsupported driver error")
	}
}
