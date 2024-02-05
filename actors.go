package actors

import (
	"log"
	"sync"
)

/*
Defines the actors and its functions
actors can receive messages and send messages
*/

// Actor ... defines a Actor type
type Actor struct {
	id     string      //actor id
	recvCh chan string //receive message on this channel
}

// ActorHub ... controls message sending among actors
type ActorHub struct {
	store     map[string]*Actor //actor store
	eventChan chan Event
}

// create a singleton isntance of the ActorHub
var actorHub *ActorHub

// used to run a finction only once
var once sync.Once

// generate a singleton ActorHub instance
func getActorsHub() (*ActorHub, error) {
	once.Do(func() {
		actorHub = &ActorHub{
			store:     make(map[string]*Actor, 100),
			eventChan: make(chan Event, 1000),
		}
		//start the actors hub once the actors instance is created
		go func() {
			actorHub.run()
		}()
	})
	if actorHub.store == nil {
		return nil, ActorError{err: ACTORHUBGENERROR}
	}
	return actorHub, nil
}

// start listening for commands
func (ah *ActorHub) run() {

	if ah == nil {
		log.Println("cannot run actors hub")
		return
	}
	for ev := range ah.eventChan {
		switch ev.eventType {
		case ADDACTOR:
			ah.registerActor(ev.actor)

		case REMOVEACTOR:
			ah.removeActor(ev.iD)

		case CLEARACTORS:
			// Delete all members of the map
			for key := range ah.store {
				delete(ah.store, key)
			}
		case SENDMESSAGE:
			ah.sendMessage(ev.iD, ev.data)

		}
	}
}

// registerActor
// if an actor already exists in the hub the regitered actors referance is returned
// / else a new actor is registered under that id
func (ah *ActorHub) registerActor(actor Actor) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	if _, ok := ah.store[actor.id]; ok {
		//an actor already exist
		return nil
	} else {
		ah.store[actor.id] = &actor
	}
	return nil
}

// AddActor ... Global function to add an actor
func AddActor(id string) (<-chan string, error) {
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
	return actorChan, nil
}

// RemoveActor ... remove an actor from the system
// de register and close the actor, so actor will not get any more messages
func (ah *ActorHub) removeActor(id string) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	//get the actor
	if actor, ok := ah.store[id]; ok {

		//close the actor channel
		close(actor.recvCh)
		//remove the actor from the stor
		delete(ah.store, id)
	}
	return nil
}

// RemoveActor ... remove an actor from the system
func RemoveActor(id string) error {
	ah, err := getActorsHub()
	if err != nil {
		return err
	}
	if ah.eventChan != nil {
		ah.eventChan <- Event{eventType: REMOVEACTOR, iD: id}
	}
	return nil
}

// SendMessage ... send message to an actor

func (ah *ActorHub) sendMessage(to string, message []byte) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	//get the actor
	if actor, ok := ah.store[to]; ok {

		if actor != nil && actor.recvCh != nil {
			actor.recvCh <- string(message)
		}

	} else {

		return ActorError{err: ACTORNOTFOUNDERROR}
	}
	return nil
}

// SendMessage ... send a message to an actor with givern id
// provide the id of the actor to which the message needs to be sent and the data in bytes
func SendMessage(id string, data []byte) error {
	ah, err := getActorsHub()
	if err != nil {
		return err
	}
	if ah.eventChan != nil {

		ah.eventChan <- Event{iD: id, data: data, eventType: SENDMESSAGE}
	}
	return nil
}
