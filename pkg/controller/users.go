package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"calendar/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type CreateUserReq struct {
	Login     string `json:"login" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type CreateUserResponse struct {
	UserId int `json:"user_id"`
}

func (c *CalendarController) CreateUser(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	req := CreateUserReq{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	err = c.Validator.StructCtx(ctx, req)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := c.UserService.CreateUser(ctx, model.User{
		Login:     req.Login,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, CreateUserResponse{UserId: id})
}

func (c *CalendarController) GetUserEvents(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	exist, err := c.UserService.IsUserExist(ctx, userID)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}
	if !exist {
		c.AbortWithBaseErrorJson(ctx, errors.Wrapf(
			model.ErrUserNotExist, "userID: %d", userID), http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, ctx.Query("from"))
	if err != nil || from.IsZero() {
		c.AbortWithBaseErrorJson(ctx, errors.Wrapf(
			model.ErrParseTimeInRequest, "err: %e, time parsed from: %s", err, from), http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, ctx.Query("to"))
	if err != nil || to.IsZero() {
		c.AbortWithBaseErrorJson(ctx, errors.Wrapf(
			model.ErrParseTimeInRequest, "err: %e, time parsed to: %s", err, to), http.StatusBadRequest)
		return
	}

	events, err := c.EventService.GetEventsByUserID(ctx, userID, from.UTC(), to.UTC())
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, events)
}
