package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"cryptoColony/src/model"
	"cryptoColony/src/storage"
)

type EventService struct {
	Repository storage.RepositoryInterface
}

func NewEventService(repository storage.RepositoryInterface) EventService {
	return EventService{Repository: repository}
}

func (e *EventService) CreateEvent(ctx context.Context, event model.Event, invitedUsers []int) (int64, error) {
	event.Duration = int(event.EndTime.Sub(event.BeginTime).Minutes())
	eventID, err := e.Repository.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
	}
	usersEvents := make([]model.UserEvent, len(invitedUsers))
	for i, v := range invitedUsers {
		usersEvents[i] = model.UserEvent{
			UserID:  v,
			EventID: eventID,
			Status:  model.NotAnswered,
		}
	}
	usersEvents = append(usersEvents, model.UserEvent{
		UserID:  event.Author,
		EventID: eventID,
		Status:  model.Accepted})
	err = e.Repository.CreateUsersEvents(ctx, usersEvents)
	return eventID, err
}

func (e *EventService) GetEvent(ctx context.Context, eventID int) (model.Event, error) {
	return e.Repository.GetEvent(ctx, eventID)
}

func (e *EventService) GetEventsByUserID(ctx context.Context, userID int) ([]model.Event, error) {
	return e.Repository.GetEventsByUserID(ctx, userID)
}

func (e *EventService) ChangeUserEventStatus(ctx context.Context, eventID, userID int, status model.InvitationStatus) error {
	userEvent, err := e.Repository.ChangeUserEventStatus(ctx, eventID, userID, status)

	if err != nil || userEvent.EventID == 0 {
		return errors.Wrap(model.ErrGettingUserEventFromDatabase, fmt.Sprintf("userID:%d, eventID:%d", userID, eventID))
	}
	if userEvent.Status != model.NotAnswered {
		return errors.Wrap(model.ErrInvitationAlreadyAnswered, fmt.Sprintf("userID:%d, eventID:%d", userID, eventID))
	}

	return nil
}
