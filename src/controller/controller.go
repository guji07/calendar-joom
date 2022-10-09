package controller

import (
	"cryptoColony/src/model"
	"cryptoColony/src/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CalendarController struct {
	UserService  service.UserService
	EventService service.EventService
	Logger       *logrus.Logger
	Validator    *validator.Validate
}

func NewCalendarController(userService service.UserService, eventService service.EventService,
	logger *logrus.Logger, validate *validator.Validate) CalendarController {
	return CalendarController{
		UserService:  userService,
		EventService: eventService,
		Logger:       logger,
		Validator:    validate,
	}
}

func (c *CalendarController) AbortWithBaseErrorJson(ctx *gin.Context, err error, status int) {
	switch {
	case errors.Is(err, model.ErrInvitationAlreadyAnswered):
		status = http.StatusConflict
	case errors.Is(err, model.ErrGettingEvent):
		status = http.StatusNotFound
	}
	c.Logger.Error(err)
	ctx.AbortWithStatusJSON(status, err.Error())
}
