package eventstore

// Writer ...
type Writer interface {
	Write(event *Event) (int, error)
}
