An actor package in Go. Use it for scalable message passing with in a process.
Actors are created under a tag and in return you will get an actor instance and you can call methods on this instance to send and receive messages over it.

Actors are assigned an internal ID and you can get it by calling actor.GetId() under the actor instance methods

Send a message to an actor by passing the tag of the actor to which you want to send the message.

You can add multiple actors for the same tag to create a room of actors and broadcast message to all the actors under that room.

You can create nodes across network by creating a master node and having other nodes joing the master node. Then you can send messages under a tag within the cluster

It is a fire and forget so messages are not saved or retried on network issues. There is no acknowledgement as well. This is a future work and is in pipeline

usage: go get github.com/michaelmenon/actors/cmd

```
import (
    github.com/michaelmenon/actors/cmd
)
//first create a Hub instance
ah := cmd.GetActorsHub()


actor1, err := ah.NewActor("actortag1")

//create another actor
actor2, err := ah.NewActor("actortag2")

//listen for messages 
for msg := <-actor1.Get(){
    fmt.Println(msg)
}
//in a seperate Goroutine send a message from actor1 to actor2 wiht the same node(same service)
actor1.SendLocal("actortag2", []byte("my message"))

//to remove an actor
actor1.Close()
actor2.Close()

//To clear all the actors 
ah.Clear()

//To creaate a cluster node
//if this node is the master then keep the MasterAddress field empty
//if you want this node to connect to a Master node provide the MasterAddress 
ah.ConnectToCluster(cmd.ClusterConfig{
		NodeName:   node,
		Host:       host,
		Port:       int(port),
		MasterAddress: "localhost:8080",
	})

We take care of the node leaving and joining with Hashicorp package that implements Gossip protocol

Just create the nodes, join a cluster by providing the MasterBode node name and send messgaes across the network

To send message to a node in a cluster under a tag
actor1.SendRemote("nodename","tagname",[]byte("data"))
```

for testing the load function run :
go test -v -run TestSendMessage

TODO : 
1) Provide reliable messaging by saving the messgaes which are not sent over the network and implement acknowledgement for received messages else retry.

2) Send messages to all nodes in the cluster


