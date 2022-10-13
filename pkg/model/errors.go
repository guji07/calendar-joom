package model

import "errors"

var (
	ErrInvitationAlreadyAnswered    = errors.New("user already answered on invitation")
	ErrGettingUserEventFromDatabase = errors.New("can't get userEvent from database")
	ErrGettingEvent                 = errors.New("no such event")
	ErrUserNotExist                 = errors.New("user does not exist")
	ErrParseTimeInRequest           = errors.New("can't correctly parse time in request")
	ErrStartTransaction             = errors.New("can't start transaction to database")
)
