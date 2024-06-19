package application

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
}

func New() *App {
	return &App{
		router: loadRoutes(),
	}
}

func (app *App) Start() error {
	err := app.router.Run(":3000")

	return err
}
