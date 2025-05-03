package model

type Screen struct {
	Text    string     `json:"text"`
	Buttons [][]Button `json:"buttons"`
}
