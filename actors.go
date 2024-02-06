package actors

const MAX_ACTORS_FOR_TAG = 100000

// Actor ... defines a Actor type
// there can be multiple actors in the system with same id but with different tag.
// tag is assigned automatically by the system and is always previous actor's tag +1
type Actor struct {
	id uint // internal id of this actor under the same id
	//can have multiple actors with different tag under the same id
	next   *Actor      //next actor with the same id but with different tag
	recvCh chan string //receive message on this channel
	tag    string      //actor tag under which it is stored
}

// AddActor ... Global function to add an actor
func NewActor(tag string) (*Actor, error) {
	ah, err := getActorsHub()
	if err != nil {
		return nil, err
	}
	actorChan := make(chan string, 1000)

	actor := Actor{
		tag:    tag,
		recvCh: actorChan,
	}
	if ah.eventChan != nil {
		ah.eventChan <- Event{tag: tag, eventType: ADDACTOR, actor: &actor}
	}
	return &actor, nil
}

// Close ... remove an actor from the system
// Close actor can only be called from an actor instance
func (a *Actor) Close() error {

	if a == nil {
		return ActorError{err: ACTORHUBGENERROR}
	}
	ah, err := getActorsHub()
	if err != nil {
		return err
	}
	if ah.eventChan != nil {

		ah.eventChan <- Event{id: a.id, tag: a.tag, eventType: REMOVEACTOR}
	}

	return nil
}

// GetReciver ... returns the receiver associated with this actor
func (a *Actor) Get() <-chan string {
	return a.recvCh
}

// GetId ... returns the id of an actor
func (a *Actor) GetTag() string {
	return a.tag
}

// GetId ... returns the id of an actor
func (a *Actor) GetId() uint {
	return a.id
}

// SendMessage ... send a message to an actor with givern id
// provide the id of the actor to which the message needs to be sent and the data in bytes
// to : the id of the actor to which we need to send the message
// data : which we need to send
func (a *Actor) SendMessage(to string, data []byte) error {
	ah, err := getActorsHub()
	if err != nil {
		return err
	}
	if ah.eventChan != nil {

		ah.eventChan <- Event{tag: to, data: data, eventType: SENDMESSAGE}
	}
	return nil
}
