package controller

import (
	"errors"
	"log"
	"net/http"

	"projectwebcurhat/config/middleware"
	"projectwebcurhat/config/pkg/errs"
	"projectwebcurhat/contract"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetPrefix() string
	InitService(service *contract.Service)
	InitRoute(app *gin.RouterGroup)
}

func New(app *gin.Engine, service *contract.Service) {
	allController := []Controller{
		&HealthController{},
		&WebSocketController{},
	}

	for _, c := range allController {
		c.InitService(service)
		group := app.Group(c.GetPrefix())
		group.Use(middleware.CORSMiddleware())
		c.InitRoute(group)
		log.Printf("initiate route %s\n", c.GetPrefix())
	}
}

func HandlerError(ctx *gin.Context, err error) {
	var messageErr errs.MessageError
	if errors.As(err, &messageErr) {
		ctx.JSON(messageErr.Status(), messageErr)
		return
	}
	_ = ctx.Error(err).SetType(gin.ErrorTypePrivate)
	ctx.JSON(http.StatusInternalServerError, errs.InternalServerError("Internal Server Error"))
}
