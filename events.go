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
	id        uint //the id of the actor
	eventType EVENTTYPE
	actor     *Actor //actor if needs to be added
	data      []byte
	tag       string //actor id to which the event is targeted
}
