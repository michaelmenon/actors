package actors

type EVENTTYPE int

const (
	ADDACTOR EVENTTYPE = iota
	REMOVEACTOR
	SENDMESSAGE
	CLEARACTORS
)

// /Event .. is an actor event.
type Event struct {
	eventType EVENTTYPE
	iD        string //actor id to which the event is targeted
	data      []byte
	actor     Actor //actor if needs to be added
}
