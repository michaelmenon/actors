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

	var err error
	if ah == nil {
		log.Println("cannot run actors hub")
		return
	}
	for ev := range ah.eventChan {
		switch ev.eventType {
		case ADDACTOR:
			err = ah.registerActor(ev.actor)
			if err != nil {
				log.Println(err)
			}

		case REMOVEACTOR:
			err = ah.removeActor(ev.iD, ev.tag)
			if err != nil {
				log.Println(err)
			}

		case CLEARACTORS:
			// Delete all members of the map
			for key := range ah.store {
				delete(ah.store, key)
			}
		case SENDMESSAGE:
			err = ah.sendMessage(ev.iD, ev.data)
			if err != nil {
				log.Println(err)
			}

		}
	}
}

// getNodeWithTag ... gets the node given a tag and an id
// returns the Actor and the previous actor
func (ah *ActorHub) getNodeWithTag(tag uint, id string) (*Actor, *Actor) {
	var prevNode *Actor
	var currNode *Actor
	if startNode, ok := ah.store[id]; ok {
		//an actor already exist
		if startNode == nil {
			return nil, nil
		}
		currNode = startNode

		for {

			if currNode.tag == tag {
				return currNode, prevNode
			}
			if currNode.next == nil {
				return nil, nil
			}
			//go to the next node, until we get the last node in the linked list
			prevNode = currNode
			currNode = currNode.next

		}
	}
	return nil, nil
}

// registerActor
// if an actor already exists in the hub the regitered actors referance is returned
// / else a new actor is registered under that id
func (ah *ActorHub) registerActor(actor *Actor) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	if oldActor, ok := ah.store[actor.id]; ok {
		//an actor already exist
		node := oldActor
		count := 0
		for {
			if node == nil {
				break
			}
			//go to the next node, until we get the last node in the linked list
			node = node.next
			count += 1
		}
		if node != nil {
			//got the last node
			//attach the new actor to this node
			//check if we have reached the max capacit of actors with the same node
			if count < MAX_ACTORS_FOR_TAG-1 {
				node.next = actor
				actor.tag = node.tag + 1
			} else {
				return ActorError{err: ACTORMAXLIMITERROR}
			}

		}
		return nil
	} else {
		//it is the root node
		actor.tag = 1
		ah.store[actor.id] = actor
	}
	return nil
}

// RemoveActor ... remove an actor from the system
// de register and close the actor, so actor will not get any more messages
func (ah *ActorHub) removeActor(id string, tag uint) error {
	if ah.store == nil {
		return ActorError{err: NILSTOREERROR}
	}
	//get the actor with the tag and id
	if actor, prevNode := ah.getNodeWithTag(tag, id); actor != nil {
		if prevNode != nil {
			prevNode.next = actor.next
		}
		//close the actor channel
		close(actor.recvCh)
		if actor.tag == 1 {
			//remove the actor from the stor
			//it is the only node
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
		for {
			if actor == nil {
				break
			}
			if actor != nil && actor.recvCh != nil {

				actor.recvCh <- string(message)
			}
			actor = actor.next
		}

	} else {

		return ActorError{err: ACTORNOTFOUNDERROR}
	}
	return nil
}
