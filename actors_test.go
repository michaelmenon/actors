package actors

import (
	"fmt"
	"testing"
	"time"
)

func TestActorCreation(t *testing.T) {

	id := "testactor"
	act, _ := NewActor(id)
	// Check if the result is not nil
	if act == nil {
		t.FailNow()
	}

}
func TestActorSendMessage(t *testing.T) {

	id := "testactor"
	act1, _ := NewActor(id)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	act2, _ := NewActor(id)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}
	err := act2.SendMessage(id, []byte("testmessage"))
	if err != nil {
		t.FailNow()
	}

}

func TestActorReceiveMessage(t *testing.T) {

	id1 := "testactor1"
	act1, _ := NewActor(id1)
	// Check if the result is not nil
	if act1 == nil {
		t.FailNow()
	}
	id2 := "testactor2"
	act2, _ := NewActor(id2)
	// Check if the result is not nil
	if act2 == nil {
		t.FailNow()
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		err := act2.SendMessage(id1, []byte("testmessage"))
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

func TestSendMessage(b *testing.T) {
	actors := make([]*Actor, 0, 100)
	//Add a million actors
	for i := 0; i < 1000000; i++ {
		id := fmt.Sprintf("%d", i)
		act1, _ := NewActor(id)
		actors = append(actors, act1)

	}
	b.Log("created all actors")
	b.Log(len(actors))
	close := "CLOSE"
	id := fmt.Sprintf("%d", 5)
	t1 := time.Now()
	go func() {

		for i, act := range actors {

			if i == 5 {
				continue
			}
			err := act.SendMessage(id, []byte(fmt.Sprintf("%d", i)))
			if err != nil {
				b.Log(err)
			}
		}
		err := actors[0].SendMessage(id, []byte(close))
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
