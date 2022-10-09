package storage

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"

	"cryptoColony/src/model"
)

func (r *Repository) CreateEvent(ctx context.Context, event model.Event) (int64, error) {
	var id int64

	_, err := r.storage.Insert(EventsTableName).Rows(event).Returning("id").Executor().ScanValContext(ctx, &id)

	return id, err
}

func (r *Repository) CreateUsersEvents(ctx context.Context, usersEvents []model.UserEvent) error {
	_, err := r.storage.Insert(UsersEventsTableName).Rows(usersEvents).Executor().ExecContext(ctx)
	return err
}

func (r *Repository) GetEvent(ctx context.Context, eventID int) (model.Event, error) {
	var event model.Event

	_, err := r.storage.Select("e.*").From(EventsTable.As("e")).Where(goqu.Ex{"id": eventID}).ScanStructContext(ctx, &event)
	if err != nil {
		return model.Event{}, err
	}
	return event, nil
}

func (r *Repository) GetEventsByUserID(ctx context.Context, userID int) ([]model.Event, error) {
	var events []model.Event
	query := r.storage.Select("e.*", "ue.status").From(EventsTable.As("e")).
		Join(UsersEventsTable.As("ue"),
			goqu.On(goqu.I("e.id").Eq(goqu.I("ue.event_id")))).
		Where(goqu.Ex{"ue.user_id": userID})
	print(query.ToSQL())

	err := query.ScanStructsContext(ctx, &events)
	if err != nil {
		return []model.Event{}, err
	}

	return events, nil
}

func (r *Repository) ChangeUserEventStatus(ctx context.Context, eventID, userID int, status model.InvitationStatus) error {
	var userEvent model.UserEvent
	haveEvent, err := r.storage.Select().From(UsersEventsTable).Where(goqu.Ex{
		"event_id": eventID,
		"user_id":  userID}).ScanStructContext(ctx, &userEvent)
	if err != nil || !haveEvent {
		return errors.Wrap(model.ErrGettingUserEventFromDatabase, fmt.Sprintf("userID:%d, eventID:%d", eventID, userID))
	}
	if userEvent.Status != model.NotAnswered {
		return errors.Wrap(model.ErrInvitationAlreadyAnswered, fmt.Sprintf("userID:%d, eventID:%d", userID, eventID))
	}
	_, err = r.storage.
		Update(UsersEventsTableName).
		Where(goqu.Ex{
			"event_id": eventID,
			"user_id":  userID,
			"status":   model.NotAnswered,
		}).Set(goqu.Ex{"status": status}).Executor().ExecContext(ctx)

	return err
}
