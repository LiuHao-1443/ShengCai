//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"shengcai/internal/handler"
	"shengcai/internal/repository"
	"shengcai/internal/server"
	"shengcai/internal/service"
	"shengcai/pkg/app"
	"shengcai/pkg/jwt"
	"shengcai/pkg/log"
	"shengcai/pkg/server/http"
	"shengcai/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewShengCaiRepository,
	repository.NewSheetInfoRepository,
	repository.NewCellDataRepository,
	repository.NewAIRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewFeiShuService,
	service.NewShengCaiService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewShengCaiHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
