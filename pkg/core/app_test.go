package core

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/config"
	"github.com/rafialfauzan/layrd/pkg/database"
)

type migrationTestModule struct {
	migrated bool
}

func (*migrationTestModule) Name() string                    { return "migration-test" }
func (*migrationTestModule) RegisterRoutes(*gin.RouterGroup) {}
func (m *migrationTestModule) AutoMigrate(*App) error {
	m.migrated = true
	return nil
}

func TestModuleMigrationsAreExplicit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := NewApp(&config.Config{GinMode: gin.TestMode, APIPrefix: "/api"})
	module := &migrationTestModule{}

	if err := app.RegisterModule(module); err != nil {
		t.Fatalf("RegisterModule() error = %v", err)
	}
	if module.migrated {
		t.Fatal("module migrated during registration")
	}
	if err := app.MigrateModules(); err != nil {
		t.Fatalf("MigrateModules() error = %v", err)
	}
	if !module.migrated {
		t.Fatal("MigrateModules() did not invoke the module migrator")
	}
}

type failingInitializerModule struct {
	err error
}

func (*failingInitializerModule) Name() string                    { return "failing-initializer" }
func (*failingInitializerModule) RegisterRoutes(*gin.RouterGroup) {}
func (m *failingInitializerModule) Init(*App) error               { return m.err }

func TestRegisterModuleReturnsInitializerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := NewApp(&config.Config{GinMode: gin.TestMode, APIPrefix: "/api"})
	wantErr := errors.New("initialization failed")

	err := app.RegisterModule(&failingInitializerModule{err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("RegisterModule() error = %v, want %v", err, wantErr)
	}
}

func TestInitDatabasesContinuesWhenOptionalDatabaseIsUnavailable(t *testing.T) {
	app := &App{
		Config: &config.Config{
			GinMode: gin.ReleaseMode,
			Databases: map[string]config.DatabaseConfig{
				"user": {
					Name:     "user",
					Optional: true,
					Driver:   "mysql",
					Host:     "127.0.0.1",
					Port:     "0",
				},
			},
		},
		DBManager: database.NewDBManager(),
	}

	if err := app.InitDatabases(); err != nil {
		t.Fatalf("InitDatabases() error = %v, want nil for unavailable optional database", err)
	}

	if _, err := app.DBManager.Get("user"); err == nil {
		t.Fatal("optional database should not be registered after connection failure")
	}
}

func TestInitDatabasesFailsWhenRequiredDatabaseIsUnavailable(t *testing.T) {
	app := &App{
		Config: &config.Config{
			GinMode: gin.ReleaseMode,
			Databases: map[string]config.DatabaseConfig{
				"default": {
					Name:   "default",
					Driver: "mysql",
					Host:   "127.0.0.1",
					Port:   "0",
				},
			},
		},
		DBManager: database.NewDBManager(),
	}

	if err := app.InitDatabases(); err == nil {
		t.Fatal("InitDatabases() error = nil, want error for unavailable required database")
	}
}
