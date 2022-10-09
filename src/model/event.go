package model

import "time"

type RepeatOptions string

type Event struct {
	Id           int64            `json:"id" db:"id" goqu:"skipinsert"`
	Name         string           `json:"name" db:"name"`
	Author       int              `json:"author" db:"author"`
	Repeatable   bool             `json:"repeatable" db:"repeatable"`
	RepeatOption RepeatOptions    `json:"repeat_options" db:"repeat_options"`
	BeginTime    time.Time        `json:"begin_time" db:"begin_time"`
	EndTime      time.Time        `json:"end_time" db:"end_time"`
	Duration     int              `json:"duration" db:"duration"`
	IsPrivate    bool             `json:"is_private" db:"is_private"`
	Details      string           `json:"details" db:"details"`
	Status       InvitationStatus `json:"status" db:"status" goqu:"skipinsert"`
}

type InvitationStatus = string

var (
	Accepted    InvitationStatus = "accepted"
	Declined    InvitationStatus = "declined"
	NotAnswered InvitationStatus = "not_answered"
)

type UserEvent struct {
	Id      int64            `db:"id" goqu:"skipinsert"`
	UserID  int              `db:"user_id"`
	EventID int64            `db:"event_id"`
	Status  InvitationStatus `db:"status"`
}
