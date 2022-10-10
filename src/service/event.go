package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teambition/rrule-go"
	"sort"
	"time"

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
	exist, err := e.Repository.IsUserExist(ctx, event.Author)
	if err != nil || !exist {
		return 0, errors.Wrapf(model.ErrUserNotExist, "userID:%d", event.Author)
	}

	event.Duration = int(event.EndTime.Sub(event.BeginTime).Minutes())
	eventID, err := e.Repository.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
	}
	usersEvents := make([]model.UserEvent, len(invitedUsers))
	for i, v := range invitedUsers {
		exist, err := e.Repository.IsUserExist(ctx, v)
		if err != nil || !exist {
			return 0, errors.Wrapf(model.ErrUserNotExist, "userID:%d", v)
		}
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
	event, err := e.Repository.GetEvent(ctx, eventID)
	if err != nil {
		return model.Event{}, err
	}
	if event.Id == 0 {
		return model.Event{}, errors.Wrap(model.ErrGettingEvent, fmt.Sprintf("eventID:%d", eventID))
	}
	return event, nil
}

func (e *EventService) GetEventsByUserID(ctx context.Context, userID int, from, to time.Time) ([]model.Event, error) {
	events, err := e.Repository.GetEventsByUserIDs(ctx, []int{userID}, &from, &to)
	if err != nil {
		return []model.Event{}, err
	}

	for i := 0; i < len(events); i++ {
		event := events[i]
		if event.IsPrivate && events[i].Author != userID {
			events[i].Details = ""
		}
		if event.Repeatable {
			ret, err := rrule.StrToRRule(event.RepeatOption)
			if err != nil {
				return []model.Event{}, err
			}
			occurrences := ret.Between(from.Add(-(time.Duration(event.Duration) * time.Minute)), to, true)
			events = RemoveEvent(events, i)
			i--
			e.addOccurrencesToEvents(occurrences, event, &events)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].BeginTime.Before(events[j].BeginTime)
	})
	return events, nil
}

func RemoveEvent(s []model.Event, index int) []model.Event {
	return append(s[:index], s[index+1:]...)
}

func (e *EventService) addOccurrencesToEvents(occurrences []time.Time, event model.Event, events *[]model.Event) {
	event.Repeatable = false
	for _, v := range occurrences {
		event.BeginTime = v.UTC()
		event.EndTime = event.BeginTime.Add(time.Duration(event.Duration) * time.Minute).UTC()
		*events = append(*events, event)
	}
}

func (e *EventService) ChangeUserEventStatus(ctx context.Context, eventID, userID int, status model.InvitationStatus) error {
	userEvent, err := e.Repository.ChangeUserEventStatus(ctx, eventID, userID, status)

	if err != nil || userEvent.EventID == 0 {
		return errors.Wrapf(model.ErrGettingUserEventFromDatabase, "userID:%d, eventID:%d, err: %v", userID, eventID, err)
	}
	if userEvent.Status != model.NotAnswered {
		return errors.Wrapf(model.ErrInvitationAlreadyAnswered, "userID:%d, eventID:%d", userID, eventID)
	}

	return nil
}
