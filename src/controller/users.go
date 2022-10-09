package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"cryptoColony/src/model"

	"github.com/gin-gonic/gin"
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
