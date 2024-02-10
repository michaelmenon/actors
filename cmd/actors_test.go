package cmd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestStartHub(t *testing.T) {
	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}
	//wait for the go rotutine to start
	time.Sleep(100 * time.Millisecond)
	tag := "testactor"
	act, _ := ah.NewActor(tag)
	// Check if the result is not nil
	if act == nil {
		t.FailNow()
	}

}

func TestActorCreation(t *testing.T) {
	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}

	tag := "testactor"
	act, _ := ah.NewActor(tag)
	// Check if the result is not nil
	if act == nil {
		t.FailNow()
	}

}
func TestActorSendMessage(t *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}
	tag := "testactor"
	act1, _ := ah.NewActor(tag)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	act2, _ := ah.NewActor(tag)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}
	err = act2.SendLocal(tag, []byte("testmessage"))
	if err != nil {
		t.FailNow()
	}

}

func TestActorReceiveMessage(t *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}

	tag1 := "testactor1"
	act1, _ := ah.NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	tag2 := "testactor2"
	act2, _ := ah.NewActor(tag2)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		err := act2.SendLocal(tag1, []byte("testmessage"))
		if err != nil {
			t.Log(err)
		}
	}()
	select {
	case <-ticker.C:
		t.Log("test failed message did not arrive in 15 seconds")
		t.FailNow()
	case msg := <-act1.Get():
		t.Logf("message recieved :%s", msg)
		if msg != "testmessage" {
			t.FailNow()
		}
	}

}

func TestBulkMessageSend(b *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		b.FailNow()
	}

	actors := make([]*Actor, 0, 100)
	//Add a million actors
	for i := 0; i < 100; i++ {
		tag := fmt.Sprintf("%d", i)
		act1, _ := ah.NewActor(tag)
		actors = append(actors, act1)

	}
	b.Log("created all actors")
	b.Log(len(actors))
	close := "CLOSE"
	tag := fmt.Sprintf("%d", 5)
	t1 := time.Now()
	go func() {

		for i, act := range actors {

			if i == 5 {
				continue
			}
			err := act.SendLocal(tag, []byte(fmt.Sprintf("%d", i)))
			if err != nil {
				b.Log(err)
			}
		}
		err := actors[0].SendLocal(tag, []byte(close))
		if err != nil {
			b.Log(err)
		}

	}()
	for {

		msg, ok := <-actors[5].Get()

		if !ok {
			break
		} else {
			b.Logf("message recieved :%s", msg)
			if msg == close {
				break
			}
		}

	}

	b.Log("Total time to receive messages:", time.Since(t1).Seconds())
}
func TestAddingNodesWithSameTag(t *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}

	tag1 := "testactor1"
	act1, _ := ah.NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}

	act2, _ := ah.NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	tag2 := "testactor2"
	act3, _ := ah.NewActor(tag2)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		err := act3.SendLocal(tag1, []byte("testmessage"))
		if err != nil {
			t.Log(err)
		}
	}()
	count := 0
	for {
		if count == 2 {
			break
		}
		select {
		case <-ticker.C:
			t.Log("test failed message did not arrive in 15 seconds")
			t.FailNow()
		case msg := <-act1.Get():
			t.Logf("message recieved for act1 :%s", msg)
			if msg != "testmessage" {
				t.FailNow()
			}
			count++
		case msg := <-act2.Get():
			t.Logf("message recieved for act2 :%s", msg)
			if msg != "testmessage" {
				t.FailNow()
			}
			count++
		}
	}

}
func TestBroadcast(b *testing.T) {
	ah, err := GetActorsHub()
	if err != nil {
		b.FailNow()
	}

	actors := make([]*Actor, 0, 100)
	tag := "tag1"
	//Add a million actors
	for i := 0; i < 100000; i++ {

		act1, _ := ah.NewActor(tag)
		actors = append(actors, act1)

	}
	b.Log("created all actors")
	b.Log(len(actors))

	wg := &sync.WaitGroup{}
	newactor, _ := ah.NewActor("tag2")
	t1 := time.Now()
	for i, act := range actors {
		var index = i
		var actor = act
		wg.Add(1)
		go func() {

			msg := <-actor.Get()
			b.Logf("got message %s for %d", msg, index)
			wg.Done()
		}()
	}

	newactor.SendLocal("tag1", []byte("test message to tag1"))
	wg.Wait()

	b.Log("Total time to receive messages:", time.Since(t1).Seconds())
}

func TestDeleteActorFromSameTag(t *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}

	actors := make([]*Actor, 0, 100)
	tag := "tag1"
	//Add a million actors
	for i := 0; i < 10; i++ {

		act1, _ := ah.NewActor(tag)
		actors = append(actors, act1)

	}
	t.Log("created all actors")
	t.Log(len(actors))
	//remove the 6th actor
	err = actors[5].Close()
	if err != nil {
		t.FailNow()
	}
	wg := &sync.WaitGroup{}
	newactor, _ := ah.NewActor("tag2")
	t1 := time.Now()
	for _, act := range actors {

		var actor = act
		wg.Add(1)
		go func() {

			msg, ok := <-actor.Get()
			if !ok {
				t.Logf("actor %d closed", actor.GetId())
			} else {
				t.Logf("got message %s for %d", msg, actor.GetId())
			}

			wg.Done()
		}()
	}

	newactor.SendLocal("tag1", []byte("test message to tag1"))
	wg.Wait()

	t.Log("Total time to receive messages:", time.Since(t1).Seconds())
}

func TestCreateCluster(t *testing.T) {

	ah, err := GetActorsHub()
	if err != nil {
		t.FailNow()
	}
	cfg := ClusterConfig{
		Host:          "127.0.0.1",
		NodeName:      "node1",
		Port:          8080,
		MasterAddress: "",
	}
	err = ah.ConnectToCluster(cfg)
	if err != nil {
		t.FailNow()
	}
}
