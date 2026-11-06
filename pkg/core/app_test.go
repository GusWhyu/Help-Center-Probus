package core

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/config"
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
