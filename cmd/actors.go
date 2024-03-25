package cmd

const MAX_ACTORS_FOR_TAG = 100000

// Actor ... defines a Actor type
// there can be multiple actors in the system with same id but with different tag.
// tag is assigned automatically by the system and is always previous actor's tag +1
type Actor struct {
	id uint // internal id of this actor under the same tag
	//can have multiple actors with different id under the same tag
	next   *Actor      //next actor with the same tag but with different id, next actor's id will have id+1
	recvCh chan string //receive message on this channel
	tag    string      //actor tag under which it is stored
}

// Close ... remove an actor from the system
// Close actor can only be called from an actor instance
func (a *Actor) Close() error {

	if a == nil {
		return ActorError{Err: ErrActorHubGen}
	}
	ah, err := GetActorsHub()
	if err != nil {
		return err
	}
	if ah.eventChan != nil {

		ah.eventChan <- Event{id: a.id, tag: a.tag, eventType: RemoveActor}
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

// SendLocal ... send a message to an actor with givern id within this node
// provide the id of the actor to which the message needs to be sent and the data in bytes
// to : the id of the actor to which we need to send the message
// data : which we need to send
func (a *Actor) SendLocal(to string, data []byte) error {
	ah, err := GetActorsHub()
	if err != nil {
		return err
	}

	if ah.ctx.Value(STATUS) == STOPPED {
		return ActorError{Err: ErrHubNotRunning}
	}
	if ah.eventChan != nil {

		ah.eventChan <- Event{tag: to, data: data, eventType: SendMessage}
	}
	return nil
}
