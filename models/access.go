package models

type AccessLevel string

const (
	AccessFreeBusy AccessLevel = "free_busy"
	AccessView     AccessLevel = "view"
)

type SharedCalendar struct {
	User  User        `json:"user"`
	Level AccessLevel `json:"level"`
}
