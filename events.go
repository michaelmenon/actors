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
	tag       string //actor id to which the event is targeted
	id        uint   //the tag of the actor to be removed
	data      []byte
	actor     *Actor //actor if needs to be added
}
