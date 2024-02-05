package actors

import (
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
