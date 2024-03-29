package storage

import (
	"context"
	"time"

	"cryptoColony/src/model"

	"github.com/doug-martin/goqu/v9"
)

var (
	UsersTableName       = "users"
	UsersTable           = goqu.T(UsersTableName)
	EventsTableName      = "events"
	EventsTable          = goqu.T(EventsTableName)
	UsersEventsTableName = "users_events"
	UsersEventsTable     = goqu.T(UsersEventsTableName)
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	IsUserExist(ctx context.Context, userID int) (bool, error)

	CreateEvent(ctx context.Context, user model.Event) (int64, error)
	GetEvent(ctx context.Context, eventID int) (model.Event, error)
	GetEventsByUserIDs(ctx context.Context, userIDs []int, from, to *time.Time) ([]model.Event, error)

	CreateUsersEvents(ctx context.Context, usersEvents []model.UserEvent) error
	ChangeUserEventStatus(ctx context.Context, eventID, userID int, status model.InvitationStatus) (model.UserEvent, error)
}

type Repository struct {
	storage *goqu.Database
}

func NewRepository(storage *goqu.Database) Repository {
	return Repository{storage: storage}
}
