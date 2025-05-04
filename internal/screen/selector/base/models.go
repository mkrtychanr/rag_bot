package selectorscreen

type SelectorOption struct {
	Option  int          `json:"option"`
	Payload SelectorData `json:"payload"`
}
