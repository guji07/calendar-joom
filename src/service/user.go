package service

import (
	"context"
	"cryptoColony/src/model"
	"cryptoColony/src/storage"
	"time"
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
	freeTimeForEvent := time.Now().Truncate(time.Minute)

	var found = false
	for !found {
		for _, userID := range usersIDs {
			isFree, err := u.CheckTimerangeForUserFree(ctx, userID, freeTimeForEvent, freeTimeForEvent.Add(time.Minute*duration))
			if err != nil {
				return time.Time{}, err
			}
			if !isFree {
				freeTimeForEvent = freeTimeForEvent.Add(time.Minute * duration)
				break
			}
		}
		found = true
	}
	return freeTimeForEvent, nil
}

func (u *UserService) CheckTimerangeForUserFree(ctx context.Context, userID int,
	startTime, endTime time.Time) (isFree bool, err error) {
	events, err := u.Repository.GetEventsByUserID(ctx, userID)

	for _, v := range events {
		if u.isTimeRangesOverlaps(startTime, endTime, v.BeginTime, v.EndTime) {
			return false, nil
		}
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserService) isTimeRangesOverlaps(startTime1, endTime1, startTime2, endTime2 time.Time) bool {
	return startTime1.Before(endTime2) && startTime2.Before(endTime1)
}
