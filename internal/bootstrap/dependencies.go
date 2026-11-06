package bootstrap

import (
	"github.com/rafialfauzan/layrd/internal/modules/user"
	"github.com/rafialfauzan/layrd/pkg/core"
)

func registerDependencies(app *core.App) {
	app.RegisterModule(user.NewModule(app))
}
