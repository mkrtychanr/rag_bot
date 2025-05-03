package model

type User struct {
	ID         int64
	TelegramID int64
	Shortname  string
}

type UserGroup struct {
	User
	AccessType RightsPolicy
}
