package actors

import (
	"log"
	"sync"
)

/*
Defines the actors Hub and its functions
actors can receive messages and send messages
*/

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
func (ah *ActorHub) registerActor(actor *Actor) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	if _, ok := ah.store[actor.id]; ok {
		//an actor already exist
		return nil
	} else {
		ah.store[actor.id] = actor
	}
	return nil
}

// RemoveActor ... remove an actor from the system
// de register and close the actor, so actor will not get any more messages
func (ah *ActorHub) removeActor(id string) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	//get the actor
	if actor, ok := ah.store[id]; ok {

		if actor != nil {
			//close the actor channel
			close(actor.recvCh)
			//remove the actor from the stor
			delete(ah.store, id)

		}

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
