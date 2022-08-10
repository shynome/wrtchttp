package sdk

type Event struct {
	ID   string
	Body []byte
}

func (e *Event) Id() string { return e.ID }

func (e *Event) Event() string { return "signaler" }

func (e *Event) Data() string { return string(e.Body) }
