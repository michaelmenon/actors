package cmd

type EventType int

const (
	AddActor EventType = iota
	RemoveActor
	SendMessage
	ClearActors
)

// /Event .. is an actor event.
type Event struct {
	id        uint //the id of the actor
	eventType EventType
	actor     *Actor //actor if needs to be added
	data      []byte
	tag       string //actor id to which the event is targeted
}
