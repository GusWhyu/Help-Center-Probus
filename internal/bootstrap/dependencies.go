package bootstrap

import (
	"github.com/rafialfauzan/layrd/internal/modules/user"
	"github.com/rafialfauzan/layrd/pkg/core"
)

func registerDependencies(app *core.App) error {
	return app.RegisterModule(user.NewModule(app))
}
