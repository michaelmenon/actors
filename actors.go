package actors

// Actor ... defines a Actor type
type Actor struct {
	id     string      //actor id
	recvCh chan string //receive message on this channel
}

// AddActor ... Global function to add an actor
func NewActor(id string) (*Actor, error) {
	ah, err := getActorsHub()
	if err != nil {
		return nil, err
	}
	actorChan := make(chan string, 10)

	actor := Actor{
		id:     id,
		recvCh: actorChan,
	}
	if ah.eventChan != nil {
		ah.eventChan <- Event{iD: id, eventType: ADDACTOR, actor: actor}
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
	ah.removeActor(a.id)
	return nil
}

// GetReciver ... returns the receiver associated with this actor
func (a *Actor) Get() <-chan string {
	return a.recvCh
}

// GetId ... returns the id of an actor
func (a *Actor) GetId() string {
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

		ah.eventChan <- Event{iD: to, data: data, eventType: SENDMESSAGE}
	}
	return nil
}
