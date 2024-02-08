An actor package in Go. Use it for scalable message passing with in a process.
Actors are created under a tag and in return you will get an actor instance and you can call methods on this instance to send and receive messages over it.

Actors are assigned an internal ID and you can get it by calling actor.GetId() under the actor instance methods

Send a message to an actor by passing the tag of the actor to which you want to send the message.

You can add multiple actors for the same tag to create a room of actors and broadcast message to all the actors under that room.

usage: go get github.com/michaelmenon/actors

```
//first create a Hub instance
ah := actors.GetActorsHub()


actor1, err := ah.NewActor("actortag1")

//create another actor
actor2, err := ah.NewActor("actortag2")

//listen for messages 
for msg := <-actor1.Get(){
    fmt.Println(msg)
}
//in a seperate Goroutine send a message from actor1 to actor2
actor1.SendMessage("actortag2", []byte("my message"))

//to remove an actor
actor1.Close()
actor2.Close()

//To clear all the actors 
ah.Clear()


```

for testing the load function run :
go test -v -run TestSendMessage



TODO : Sending messages to actors over the network


