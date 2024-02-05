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

	id1 := "testactor1"
	act1, _ := NewActor(id1)
	// Check if the result is not nil
	if act1 == nil {
		b.FailNow()
	}
	id2 := "testactor2"
	act2, _ := NewActor(id2)
	// Check if the result is not nil
	if act2 == nil {
		b.FailNow()
	}
	close := "CLOSE"
	go func() {

		for i := 0; i < 1000000; i++ {

			err := act2.SendMessage(id1, []byte(fmt.Sprintf("%d", i)))
			if err != nil {
				b.Log(err)
			}
		}
		err := act2.SendMessage(id1, []byte(close))
		if err != nil {
			b.Log(err)
		}

	}()
	for {

		msg, ok := <-act1.recvCh

		if !ok {
			break
		} else {
			b.Logf("message recieved :%s", msg)
			if msg == close {
				break
			}
		}

	}

}
