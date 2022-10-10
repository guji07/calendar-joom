package storage

import (
	"context"
	"time"

	"cryptoColony/src/model"

	"github.com/doug-martin/goqu/v9"
)

func (r *Repository) CreateEvent(ctx context.Context, event model.Event) (int64, error) {
	var id int64

	_, err := r.storage.Insert(EventsTable).Rows(event).Returning("id").Executor().ScanValContext(ctx, &id)

	return id, err
}

func (r *Repository) CreateUsersEvents(ctx context.Context, usersEvents []model.UserEvent) error {
	_, err := r.storage.Insert(UsersEventsTable).Rows(usersEvents).Executor().ExecContext(ctx)
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

func (r *Repository) GetEventsByUserIDs(ctx context.Context, userIDs []int, from, to time.Time) ([]model.Event, error) {
	var events []model.Event
	query := r.storage.
		Select("e.*", "ue.status").
		From(EventsTable.As("e")).
		Join(UsersEventsTable.As("ue"),
			goqu.On(goqu.I("e.id").Eq(goqu.I("ue.event_id")))).
		Where(
			goqu.I("ue.user_id").In(userIDs),
			goqu.I("ue.status").In(model.Accepted, model.NotAnswered),
		)

	if !from.IsZero() {
		query.Where(goqu.Or(goqu.I("e.begin_time").Lte(to), goqu.I("e.repeatable").Eq(true)))
	}
	if !to.IsZero() {
		query.Where(goqu.Or(goqu.I("e.end_time").Gte(from), goqu.I("e.repeatable").Eq(true)))
	}

	err := query.ScanStructsContext(ctx, &events)
	if err != nil {
		return []model.Event{}, err
	}

	return events, nil
}

func (r *Repository) ChangeUserEventStatus(ctx context.Context, eventID, userID int, status model.InvitationStatus) (model.UserEvent, error) {
	var userEvent model.UserEvent

	haveEvent, err := r.storage.Select().From(UsersEventsTable).Where(goqu.Ex{
		"event_id": eventID,
		"user_id":  userID}).ScanStructContext(ctx, &userEvent)
	if err != nil || !haveEvent {
		return userEvent, err
	}

	_, err = r.storage.
		Update(UsersEventsTable).
		Where(goqu.Ex{
			"event_id": eventID,
			"user_id":  userID,
			"status":   model.NotAnswered,
		}).Set(goqu.Ex{"status": status}).Executor().ExecContext(ctx)

	return userEvent, err
}
