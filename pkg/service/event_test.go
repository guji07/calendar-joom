package service

import (
	"context"
	"testing"
	"time"

	"calendar/pkg/model"
	"calendar/pkg/storage/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testUsers = []model.User{
		{
			Login:     "testLogin1",
			FirstName: "firstName",
			LastName:  "lastName",
		}, {
			Login:     "testLogin2",
			FirstName: "firstName",
			LastName:  "lastName",
		}, {
			Login:     "testLogin3",
			FirstName: "firstName",
			LastName:  "lastName",
		},
		{
			Login:     "testLogin4",
			FirstName: "firstName",
			LastName:  "lastName",
		},
		{
			Login:     "testLogin5",
			FirstName: "firstName",
			LastName:  "lastName",
		},
		{
			Login:     "testLogin6",
			FirstName: "firstName",
			LastName:  "lastName",
		},
		{
			Login:     "testLogin7",
			FirstName: "authorName",
			LastName:  "lastName",
		}}
	timeNow = time.Now()

	testEventSuccess = model.Event{
		Name:         "name",
		Author:       6,
		Repeatable:   false,
		RepeatOption: "",
		BeginTime:    timeNow,
		EndTime:      timeNow.Add(time.Hour),
		Duration:     60,
		IsPrivate:    false,
		Details:      "details",
	}
	testEventWithRepeatOptions = model.Event{
		Name:         "name",
		Author:       6,
		Repeatable:   true,
		RepeatOption: "DTSTART;TZID=Europe/Moscow:20221009T220000\\nFREQ=HOURLY;INTERVAL=1;COUNT=20",
		BeginTime:    timeNow,
		EndTime:      timeNow.Add(time.Hour),
		Duration:     60,
		IsPrivate:    false,
		Details:      "details",
	}

	testEventFailed = model.Event{
		Name:         "name",
		Author:       7,
		Repeatable:   false,
		RepeatOption: "",
		BeginTime:    timeNow,
		EndTime:      timeNow.Add(time.Hour),
		Duration:     60,
		IsPrivate:    false,
		Details:      "details",
	}

	invitedUsers    = []int{0, 1, 2, 3, 4, 5}
	testUsersEvents = []model.UserEvent{{
		UserID:  0,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  1,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  2,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  3,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  4,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  5,
		EventID: 1,
		Status:  model.NotAnswered,
	}, {
		UserID:  6,
		EventID: 1,
		Status:  model.Accepted,
	},
	}
)

func TestEventService(t *testing.T) {
	ctx := context.Background()
	repositoryMock := generateMock(ctx)
	t.Run("createEvent success", func(t *testing.T) {
		repositoryMock.On("CreateEvent", ctx, testEventSuccess).Return(int64(1), nil)
		repositoryMock.On("CreateUsersEvents", ctx, testUsersEvents).Return(nil)
		eventService := NewEventService(repositoryMock)
		userService := NewUserService(repositoryMock)
		for i, v := range testUsers {
			userID, err := userService.CreateUser(ctx, v)
			assert.Equal(t, i, userID)
			assert.NoError(t, err, "error in creating user")
		}
		eventID, err := eventService.CreateEvent(ctx, testEventSuccess, invitedUsers)
		assert.NoError(t, err, "create event error")
		assert.Equal(t, int64(1), eventID, "eventID not as expected")
	})
	repositoryMock = generateMock(ctx)
	t.Run("createEvent fail", func(t *testing.T) {
		repositoryMock.On("CreateEvent", ctx, testEventSuccess).Return(int64(1), nil)
		repositoryMock.On("IsUserExist", ctx, 7).Return(false, nil)
		eventService := NewEventService(repositoryMock)
		userService := NewUserService(repositoryMock)
		for i, v := range testUsers {
			userID, err := userService.CreateUser(ctx, v)
			assert.Equal(t, i, userID)
			assert.NoError(t, err, "error in creating user")
		}
		_, err := eventService.CreateEvent(ctx, testEventFailed, invitedUsers)
		assert.Error(t, errors.Wrapf(model.ErrUserNotExist, "userID:%d", testEventFailed.Author), err)
	})
}

func generateMock(ctx context.Context) *mocks.RepositoryInterface {
	repositoryMock := mocks.RepositoryInterface{}
	repositoryMock.On("CreateUser", mock.Anything, testUsers[0]).Return(0, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[1]).Return(1, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[2]).Return(2, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[3]).Return(3, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[4]).Return(4, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[5]).Return(5, nil)
	repositoryMock.On("CreateUser", mock.Anything, testUsers[6]).Return(6, nil)
	repositoryMock.On("IsUserExist", ctx, 6).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 0).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 1).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 2).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 3).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 4).Return(true, nil)
	repositoryMock.On("IsUserExist", ctx, 5).Return(true, nil)
	return &repositoryMock
}
