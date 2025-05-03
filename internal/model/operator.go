package model

type OperatorData struct {
	ChatID       int64
	UserID       int64
	MessageID    int64
	Text         *string
	CallbackData []byte
	DocumentID   string
}
