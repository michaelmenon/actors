package cmd

import (
	"context"
	"log"
	"log/slog"
	"sync"
)

const ActorsHub = "HUB"

// define the HUBSTATUS wheather running or stopped
type HubStatus int

const STATUS = "hubstatus"

const (
	STOPPED HubStatus = iota
	RUNNING
)

/*
Defines the actors Hub and its functions
actors can receive messages and send messages
*/

// ActorHub ... controls message sending among actors
type ActorHub struct {
	store     map[string]*Actor //actor store
	eventChan chan Event        //listen for commands from the actor
	ctx       context.Context   //context to listen for start and stop
	nodeName  string            //node name if its a part of the cluster
}

// create a singleton isntance of the ActorHub
var actorHub *ActorHub

// used to run a finction only once
var once sync.Once

func init() {
	//start the HUB while package is initialized
	GetActorsHub()

}

// generate a singleton ActorHub instance
func GetActorsHub() (*ActorHub, error) {

	once.Do(func() {
		ctx := context.Background()
		actorHub = &ActorHub{
			store:     make(map[string]*Actor, 100),
			eventChan: make(chan Event, 1000),
			ctx:       context.WithValue(ctx, STATUS, STOPPED),
		}
		go func() {
			actorHub.run()
		}()
	})
	if actorHub == nil || actorHub.store == nil {
		return nil, ActorError{Err: ErrActorHubGen}
	}
	return actorHub, nil
}

// NewActor ...  function to add an actor
func (ah *ActorHub) NewActor(tag string) (*Actor, error) {

	actorChan := make(chan string, 1000)

	actor := Actor{
		tag:    tag,
		recvCh: actorChan,
	}
	if ah.eventChan != nil {
		ah.eventChan <- Event{tag: tag, eventType: AddActor, actor: &actor}
	}
	return &actor, nil
}

// Clear clear all the data

func (ah *ActorHub) Clear() error {

	if ah == nil {
		return ActorError{Err: ErrNilStore}
	}
	ah.eventChan <- Event{eventType: ClearActors}
	return nil
}

// start listening for commands
func (ah *ActorHub) run() error {

	var err error
	if ah == nil {
		log.Println("cannot run actors hub")
		return ActorError{Err: ErrNilStore}
	}
	if ah.ctx.Value(STATUS) == RUNNING {
		//already running
		return ActorError{Err: ErrHubRunning}
	}

	ah.ctx = context.WithValue(context.Background(), STATUS, RUNNING)

	slog.Info("HUB Started")
	for ev := range ah.eventChan {
		switch ev.eventType {
		case AddActor:
			err = ah.registerActor(ev.actor)
			if err != nil {
				slog.Error(err.Error(), "event", ev)
			}

		case RemoveActor:
			err = ah.removeActor(ev.tag, ev.id)
			if err != nil {
				slog.Error(err.Error(), "event", ev)
			}

		case ClearActors:
			// Delete all members of the map
			for key := range ah.store {
				delete(ah.store, key)
			}

		case SendMessage:
			err = ah.sendMessage(ev.tag, ev.data)
			if err != nil {
				slog.Error(err.Error(), "tag", ev.tag)
			}

		}
	}
	return nil
}

// getNodeWithTag ... gets the node given a tag and an id
// returns the Actor and the previous actor
func (ah *ActorHub) getNodeWithTag(tag string, id uint) (*Actor, *Actor) {
	var prevNode *Actor
	var currNode *Actor
	if startNode, ok := ah.store[tag]; ok {
		//an actor already exist

		currNode = startNode
		//check if the first node is the one we are looking for

		for currNode != nil {

			if currNode.id == id {
				return currNode, prevNode
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
		return ActorError{Err: ErrNilStore}
	}
	//check if the hub us running or not
	if ah.ctx.Value(STATUS) == STOPPED {
		return ActorError{Err: ErrHubNotRunning}
	}
	if oldActor, ok := ah.store[actor.tag]; ok {
		//an actor already exist
		node := oldActor
		count := 0
		for node != nil && node.next != nil {

			node = node.next
			count += 1
		}
		if node != nil {
			//got the last node
			//attach the new actor to this node
			//check if we have reached the max capacit of actors with the same node
			if count < MAX_ACTORS_FOR_TAG-1 {
				node.next = actor
				actor.id = node.id + 1
			} else {
				return ActorError{Err: ErrActorMaxLimit}
			}

		}
		return nil
	} else {
		//it is the root node
		actor.id = 1
		ah.store[actor.tag] = actor
	}
	return nil
}

// RemoveActor ... remove an actor from the system
// de register and close the actor, so actor will not get any more messages
func (ah *ActorHub) removeActor(tag string, id uint) error {
	if ah.store == nil {
		return ActorError{Err: ErrNilStore}
	}
	//check if the hub us running or not
	if ah.ctx.Value(STATUS) == STOPPED {
		return ActorError{Err: ErrHubNotRunning}
	}
	//get the actor with the tag and id
	if actor, prevNode := ah.getNodeWithTag(tag, id); actor != nil {
		if prevNode != nil {
			prevNode.next = actor.next
		}
		//close the actor channel
		close(actor.recvCh)
		if actor.id == 1 {
			//remove the actor from the stor
			//it is the only node
			delete(ah.store, tag)
		}

	}
	return nil
}

// SendMessage ... send message to an actor
func (ah *ActorHub) sendMessage(to string, message []byte) error {
	if ah.store == nil {
		return ActorError{Err: ErrNilStore}
	}
	//check if the hub us running or not
	if ah.ctx.Value(STATUS) == STOPPED {
		return ActorError{Err: ErrHubNotRunning}
	}
	//get the actor
	if actor, ok := ah.store[to]; ok {
		for actor != nil {
			if actor.recvCh != nil {
				actor.recvCh <- string(message)
			}
			actor = actor.next
		}

	} else {

		return ActorError{Err: ErrActorNotFound}
	}
	return nil
}
