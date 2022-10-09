package service

import (
	"context"
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

func (u *UserService) FindWindowForUsers(ctx context.Context, usersIDs []int, duration time.Duration) (time.Time, error) {
	freeTimeForEvent := time.Now().Truncate(time.Minute).UTC()

	var found = false
	for !found {
		found = true
		for _, userID := range usersIDs {
			isFree, overlapEnd, err := u.CheckTimerangeForUserFree(ctx, userID, freeTimeForEvent, freeTimeForEvent.Add(time.Minute*duration))
			if err != nil {
				return time.Time{}, err
			}
			if !isFree {
				freeTimeForEvent = overlapEnd
				found = false
				break
			}
		}
	}
	return freeTimeForEvent, nil
}

func (u *UserService) CheckTimerangeForUserFree(ctx context.Context, userID int,
	startTime, endTime time.Time) (isFree bool, overlapEnd time.Time, err error) {
	events, err := u.Repository.GetEventsByUserID(ctx, userID)

	for _, v := range events {
		if u.isTimeRangesOverlaps(startTime, endTime, v.BeginTime.UTC(), v.EndTime.UTC()) {
			return false, v.EndTime.UTC(), nil
		}
	}
	if err != nil {
		return false, time.Time{}, err
	}
	return true, time.Time{}, nil
}

func (u *UserService) isTimeRangesOverlaps(windowStart, windowEnd, eventStart, eventEnd time.Time) bool {
	println(windowStart.String(), windowEnd.String(), eventStart.String(), eventEnd.String())
	println("check results: ", windowStart.Before(eventEnd) && eventStart.Before(windowEnd))
	return windowStart.Before(eventEnd) && eventStart.Before(windowEnd)
}
