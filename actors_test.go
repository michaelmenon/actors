package actors

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestActorCreation(t *testing.T) {

	tag := "testactor"
	act, _ := NewActor(tag)
	// Check if the result is not nil
	if act == nil {
		t.FailNow()
	}

}
func TestActorSendMessage(t *testing.T) {

	tag := "testactor"
	act1, _ := NewActor(tag)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	act2, _ := NewActor(tag)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}
	err := act2.SendMessage(tag, []byte("testmessage"))
	if err != nil {
		t.FailNow()
	}

}

func TestActorReceiveMessage(t *testing.T) {

	tag1 := "testactor1"
	act1, _ := NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	tag2 := "testactor2"
	act2, _ := NewActor(tag2)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		err := act2.SendMessage(tag1, []byte("testmessage"))
		if err != nil {
			t.Log(err)
		}
	}()
	select {
	case <-ticker.C:
		t.Log("test failed message did not arrive in 15 seconds")
		t.FailNow()
	case msg := <-act1.recvCh:
		t.Logf("message recieved :%s", msg)
		if msg != "testmessage" {
			t.FailNow()
		}
	}

}

func TestBulkMessageSend(b *testing.T) {
	actors := make([]*Actor, 0, 100)
	//Add a million actors
	for i := 0; i < 100; i++ {
		tag := fmt.Sprintf("%d", i)
		act1, _ := NewActor(tag)
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
			err := act.SendMessage(tag, []byte(fmt.Sprintf("%d", i)))
			if err != nil {
				b.Log(err)
			}
		}
		err := actors[0].SendMessage(tag, []byte(close))
		if err != nil {
			b.Log(err)
		}

	}()
	for {

		msg, ok := <-actors[5].recvCh

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

	tag1 := "testactor1"
	act1, _ := NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}

	act2, _ := NewActor(tag1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	tag2 := "testactor2"
	act3, _ := NewActor(tag2)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		err := act3.SendMessage(tag1, []byte("testmessage"))
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
		case msg := <-act1.recvCh:
			t.Logf("message recieved for act1 :%s", msg)
			if msg != "testmessage" {
				t.FailNow()
			}
			count++
		case msg := <-act2.recvCh:
			t.Logf("message recieved for act2 :%s", msg)
			if msg != "testmessage" {
				t.FailNow()
			}
			count++
		}
	}

}
func TestBroadcast(b *testing.T) {
	actors := make([]*Actor, 0, 100)
	tag := "tag1"
	//Add a million actors
	for i := 0; i < 100000; i++ {

		act1, _ := NewActor(tag)
		actors = append(actors, act1)

	}
	b.Log("created all actors")
	b.Log(len(actors))

	wg := &sync.WaitGroup{}
	newactor, _ := NewActor("tag2")
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

	newactor.SendMessage("tag1", []byte("test message to tag1"))
	wg.Wait()

	b.Log("Total time to receive messages:", time.Since(t1).Seconds())
}

func TestDeleteActorFromSameTag(t *testing.T) {

	actors := make([]*Actor, 0, 100)
	tag := "tag1"
	//Add a million actors
	for i := 0; i < 10; i++ {

		act1, _ := NewActor(tag)
		actors = append(actors, act1)

	}
	t.Log("created all actors")
	t.Log(len(actors))
	//remove the 6th actor
	err := actors[5].Close()
	if err != nil {
		t.FailNow()
	}
	wg := &sync.WaitGroup{}
	newactor, _ := NewActor("tag2")
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

	newactor.SendMessage("tag1", []byte("test message to tag1"))
	wg.Wait()

	t.Log("Total time to receive messages:", time.Since(t1).Seconds())
}
