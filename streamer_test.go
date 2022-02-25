package streamer_test

import (
	"testing"
	"time"

	"github.com/kevinjoiner/streamer"
)

func TestStream_ZeroValue(t *testing.T) {
	// defer func() {
	// 	r := recover()
	// 	if r != nil {
	// 		t.Fatalf("Paniced while using a zero valued stream: %v", r)
	// 	}
	// }()
	var zeroStream streamer.Stream
	if !zeroStream.IsZero() {
		t.Fatal("Zero intialized stream was not detected.")
	}
	if !zeroStream.IsStopped() {
		t.Fatal("Zero intialized stream was not stopped.")
	}
	zeroStream.Send("r")
	_ = zeroStream.C()
	wasStopped := zeroStream.Stop()
	if wasStopped {
		t.Fatal("Return of stop on a zero intialized stream was not true.")
	}
}

func TestStream_ZeroValueSend(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("Paniced while sending on a zero valued stream: %v", r)
		}
	}()
	var zeroStream streamer.Stream
	zeroStream.Send("r")
}

func TestStream_StoppedSend(t *testing.T) {
	// defer func() {
	// 	r := recover()
	// 	if r != nil {
	// 		t.Fatalf("Paniced while sending on a stopped stream: %v", r)
	// 	}
	// }()
	var stream streamer.Stream
	stream.Start(0)
	_ = stream.Stop()
	stream.Send("r")
	stream.C()
	stream.Start(32)
	_ = stream.Stop()
	stream.Send("r")
	stream.C()
}

func TestStream_Send(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("Paniced while sending on a stream: %v", r)
		}
	}()
	stream := streamer.New(2)
	stream.Send("foo")
	stream.Send("bar")
	stream.Send("car")
	stream.Start(0)
	stream.Send("tar")
	wasStopped := stream.Stop()
	if !wasStopped {
		t.Fatal("Return of stop on a started stream was not true.")
	}
}

func TestStream_SendBlocking(t *testing.T) {

	stream := streamer.New(2)
	stream.SendBlocking("tar")
	stream.SendBlocking("foo")
	sendChan := make(chan struct{})

	go func(c chan struct{}) {
		stream.SendBlocking("bar")
		close(c)
	}(sendChan)

	select {
	case <-sendChan:
		t.Fatal("Send blocking did not block")
	case <-time.After(time.Millisecond * 50):
	}
	wasStopped := stream.Stop()
	if !wasStopped {
		t.Fatal("Return of stop on a started stream was not true.")
	}

	sendChan = make(chan struct{})
	go func(c chan struct{}) {
		stream.SendBlocking("bar2")
		close(c)
	}(sendChan)

	select {
	case <-sendChan:
	case <-time.After(time.Millisecond * 50):
		t.Fatal("Send blocking blocked on a stopped stream")
	}
}

func TestStream_Receive(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("Paniced while sending on a stream: %v", r)
		}
	}()
	var stream streamer.Stream
	const foo = "foo"
	const bar = 85
	stream.Start(2)
	stream.Send("discarded value")
	stream.Send(foo)
	stream.Send(bar)
	value := <-stream.C()
	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("Received type %T from stream wanted %T", value, foo)
	}
	if strValue != foo {
		t.Fatalf("Received %v from stream wanted %v", value, foo)
	}

	value = <-stream.C()
	intValue, ok := value.(int)
	if !ok {
		t.Fatalf("Received type %T from stream wanted %T", value, bar)
	}
	if intValue != bar {
		t.Fatalf("Received %v from stream wanted %v", value, bar)
	}
	wasStopped := stream.Stop()
	if !wasStopped {
		t.Fatal("Return of stop on a started stream was not true.")
	}
	wasStopped = stream.Stop()
	if wasStopped {
		t.Fatal("Second call to stop returned true.")
	}
}
