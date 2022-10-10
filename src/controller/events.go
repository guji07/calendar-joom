package controller

import (
	"cryptoColony/src/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateEventReq struct {
	Name          string    `json:"name" validate:"required"`
	Author        int       `json:"author" validate:"required"`
	Repeatable    bool      `json:"repeatable"`
	RepeatOptions string    `json:"repeat_options"`
	BeginTime     time.Time `json:"begin_time" validate:"required"`
	EndTime       time.Time `json:"end_time" validate:"required,gtefield=BeginTime"`
	IsPrivate     bool      `json:"is_private"`
	Details       string    `json:"details"`
	InvitedUsers  []int     `json:"invited_users"`
}

type GetEventResp struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Author     int       `json:"author"`
	Repeatable bool      `json:"repeatable"`
	BeginTime  time.Time `json:"begin_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   int       `json:"duration"`
	IsPrivate  bool      `json:"is_private"`
	Details    string    `json:"details"`
}

type RespondOnEventReq struct{}

type FindWindowForEventReq struct {
	UsersIDs []int `json:"users_ids"`
	Duration int   `json:"duration" validator:"required"`
}

type GetUserEventsResp struct {
	Events model.Event `json:"events"`
}

type CreateEventResp struct {
	EventId int64 `json:"event_id"`
}

type FindWindowForEventResp struct {
	StartTime time.Time `json:"start_time"`
}

func (c *CalendarController) CreateEvent(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	req := CreateEventReq{}
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

	id, err := c.EventService.CreateEvent(ctx, model.Event{
		Name:         req.Name,
		Author:       req.Author,
		Repeatable:   req.Repeatable,
		RepeatOption: req.RepeatOptions,
		BeginTime:    req.BeginTime.Truncate(time.Minute),
		EndTime:      req.EndTime.Truncate(time.Minute),
		IsPrivate:    req.IsPrivate,
		Details:      req.Details,
	}, req.InvitedUsers)

	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, CreateEventResp{EventId: id})
}

func (c *CalendarController) GetEvent(ctx *gin.Context) {
	id := ctx.Param("id")

	eventID, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	if eventID == 0 {
		c.AbortWithBaseErrorJson(ctx, fmt.Errorf("eventID must be > 0"), http.StatusBadRequest)
		return
	}

	event, err := c.EventService.GetEvent(ctx, eventID)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, GetEventResp{
		Id:         event.Id,
		Name:       event.Name,
		Author:     event.Author,
		Repeatable: event.Repeatable,
		BeginTime:  event.BeginTime.UTC(),
		EndTime:    event.EndTime.UTC(),
		Duration:   event.Duration,
		IsPrivate:  event.IsPrivate,
		Details:    event.Details,
	})
}

func (c *CalendarController) RespondOnEvent(ctx *gin.Context) {
	var status model.InvitationStatus

	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
	}

	accept, err := strconv.ParseBool(ctx.Query("accept"))
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
	}

	eventID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	if accept {
		status = model.Accepted
	} else {
		status = model.Declined
	}
	err = c.EventService.ChangeUserEventStatus(ctx, eventID, userID, status)
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *CalendarController) FindWindowForEvent(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusBadRequest)
		return
	}

	req := FindWindowForEventReq{}
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

	startTime, err := c.UserService.FindWindowForUsers(ctx, req.UsersIDs, time.Duration(req.Duration))
	if err != nil {
		c.AbortWithBaseErrorJson(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, FindWindowForEventResp{startTime})
}
