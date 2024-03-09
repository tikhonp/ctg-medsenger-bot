package ctg_medsenger_bot

import (
	"fmt"
	"github.com/TikhonP/ctg-medsenger-bot/appconfig"
	"github.com/TikhonP/ctg-medsenger-bot/handler"
	"github.com/TikhonP/ctg-medsenger-bot/util"
	"github.com/TikhonP/maigo"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handlers struct {
	root     *handler.RootHandler
	init     *handler.InitHandler
	status   *handler.StatusHandler
	remove   *handler.RemoveHandler
	settings *handler.SettingsHandler
}

func createHandlers(MAIClient *maigo.Client) *handlers {
	return &handlers{
		init: &handler.InitHandler{MAIClient: MAIClient},
	}
}

func Serve(cfg *appconfig.Server) {
	MAIClient := maigo.Init(cfg.MedsengerAgentKey)
	handlers := createHandlers(MAIClient)

	app := echo.New()
	app.Debug = cfg.Debug
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.BodyDump(func(context echo.Context, req []byte, res []byte) {
		fmt.Println()
		fmt.Println("Request:", string(req))
		fmt.Println("Response:", string(res))
	}))
	if !cfg.Debug {
		app.Use(sentryecho.New(sentryecho.Options{}))
	}
	app.Validator = util.NewDefaultValidator()

	app.GET("/", handlers.root.Handle)
	app.POST("/init", handlers.init.Handle, util.ApiKeyJSON(cfg))
	app.POST("/status", handlers.status.Handle, util.ApiKeyJSON(cfg))
	app.POST("/remove", handlers.remove.Handle, util.ApiKeyJSON(cfg))
	app.GET("/settings", handlers.settings.Handle, util.ApiKeyGetParam(cfg))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	err := app.Start(addr)
	if err != nil {
		panic(err)
	}
}
