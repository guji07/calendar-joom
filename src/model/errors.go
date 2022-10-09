package model

import "errors"

var (
	ErrInvitationAlreadyAnswered    = errors.New("user already answered on invitation")
	ErrGettingUserEventFromDatabase = errors.New("can't get userEvent from database")
	ErrGettingEvent                 = errors.New("no such event")
)
