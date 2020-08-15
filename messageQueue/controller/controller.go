package controller

import (
	"context"
	"fmt"
	"messageQueue/models"
	"messageQueue/service"
	"net/http"

	"github.com/labstack/echo"
)

type controller struct {
	qS service.QueueService
}

func ApplyController(router *echo.Echo, queueService service.QueueService) {
	handler := &controller{queueService}

	check := router.Group("check")
	check.POST("", handler.ConfigureAuth)

	messageQueue := router.Group("mq")
	messageQueue.POST("", handler.pushMessage)
	messageQueue.GET("", handler.GetMessage)
	messageQueue.DELETE("", handler.DeleteQueue)
}

var auth bool

func (a *controller) pushMessage(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if auth == true {
		username, password, ok := c.Request().BasicAuth()
		fmt.Print(ok)

		if username != "admin" && password != "password" {
			return echo.ErrUnauthorized
		}

	}

	message := models.MessageIn{}
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(c.QueryParam("key"))
	if err := a.qS.PushMessage(ctx, c.QueryParam("key"), message); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, message)
}

func (a *controller) GetMessage(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if auth == true {
		username, password, ok := c.Request().BasicAuth()
		fmt.Print(ok)

		if username != "admin" && password != "password" {
			return echo.ErrUnauthorized
		}

	}

	message, err := a.qS.GetMessage(ctx, c.QueryParam("key"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, message)
}

func (a *controller) DeleteQueue(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if auth == true {
		username, password, ok := c.Request().BasicAuth()
		fmt.Print(ok)

		if username != "admin" && password != "password" {
			return echo.ErrUnauthorized
		}

	}

	err := a.qS.DeleteQueue(ctx, c.QueryParam("key"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}

func (a *controller) ConfigureAuth(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	check := models.Configure{}
	if err := c.Bind(&check); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	auth = check.IsAuth

	return c.JSON(http.StatusOK, "OK")
}
