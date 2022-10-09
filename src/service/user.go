package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/teambition/rrule-go"
	"sort"
	"time"

	"cryptoColony/src/model"
	"cryptoColony/src/storage"
)

type UserService struct {
	Repository storage.RepositoryInterface
}

func NewUserService(repository storage.RepositoryInterface) UserService {
	return UserService{Repository: repository}
}

func (u *UserService) CreateUser(ctx context.Context, user model.User) (int, error) {
	return u.Repository.CreateUser(ctx, user)
}

func (u *UserService) IsUserExist(ctx context.Context, userID int) (bool, error) {
	return u.Repository.IsUserExist(ctx, userID)
}

func (u *UserService) FindWindowForUsers(ctx context.Context, usersIDs []int, duration time.Duration) (time.Time, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	return u.searchForFreeTime(ctxWithTimeout, usersIDs, duration)
}

func (u *UserService) searchForFreeTime(ctx context.Context, usersIDs []int, duration time.Duration) (time.Time, error) {
	freeTimeForEvent := time.Now().Truncate(time.Minute).UTC()

	for {
		select {
		case <-ctx.Done():
			return time.Time{}, errors.New("context timeout in searchForFreeTime")
		default:
			isFree, overlapEnd, err := u.checkTimerangeIsFree(ctx, usersIDs, freeTimeForEvent, freeTimeForEvent.Add(time.Minute*duration))
			if err != nil {
				return time.Time{}, err
			}
			if !isFree {
				freeTimeForEvent = overlapEnd
			} else {
				return freeTimeForEvent, nil
			}
		}
	}
}

func (u *UserService) checkTimerangeIsFree(ctx context.Context, userIDs []int,
	startTime, endTime time.Time) (isFree bool, overlapEnd time.Time, err error) {
	events, err := u.Repository.GetEventsByUserIDs(ctx, userIDs, time.Time{}, startTime)
	for _, event := range events {
		if event.Repeatable {
			isFree, overlapEnd, err = u.checkRepeatableEvent(event, startTime, endTime)
			if err != nil || !isFree {
				return isFree, overlapEnd, err
			}
		} else {
			if u.isTimeRangesOverlaps(startTime, endTime, event.BeginTime.UTC(), event.EndTime.UTC()) {
				return false, event.EndTime.UTC(), nil
			}
		}
	}
	if err != nil {
		return false, time.Time{}, err
	}
	return true, time.Time{}, nil
}

func (u *UserService) checkRepeatableEvent(event model.Event, startTime time.Time, endTime time.Time) (isFree bool, overlapEnd time.Time, err error) {
	ret, err := rrule.StrToRRule(event.RepeatOption)
	if err != nil {
		return false, time.Time{}, err
	}

	occurrences := ret.Between(startTime.Add(-(time.Duration(event.Duration) * time.Minute)), endTime, false)

	if len(occurrences) > 0 {
		sort.Slice(occurrences, func(i, j int) bool { return occurrences[i].Before(occurrences[j]) })
		return false, occurrences[len(occurrences)-1].Add(time.Duration(event.Duration) * time.Minute).UTC(), nil
	}
	return true, time.Time{}, nil
}

func (u *UserService) isTimeRangesOverlaps(windowStart, windowEnd, eventStart, eventEnd time.Time) bool {
	return windowStart.Before(eventEnd) && eventStart.Before(windowEnd)
}
