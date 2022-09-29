package domain

type Event struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type EventHolder struct {
	Event *Event
	Err   error
}
