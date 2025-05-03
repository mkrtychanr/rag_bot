package model

type MenuOption struct {
	Option  int64  `json:"option"`
	Payload []byte `json:"payload"`
}
