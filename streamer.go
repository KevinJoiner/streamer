// Package streamer is used to stream data utilizing channels
package streamer

import (
	"sync"
	"time"
)

// Stream is used for sending and receiving data concurrently.
type Stream struct {
	c        chan interface{}
	stopchan chan struct{}
	l        sync.Mutex
	stoplock sync.Mutex
	started  bool
}

// New returns a new stream object with the specified buffer size.
func New(bufferSize int) *Stream {
	return &Stream{
		c:        make(chan interface{}, bufferSize),
		stopchan: make(chan struct{}),
		l:        sync.Mutex{},
		stoplock: sync.Mutex{},
		started:  true,
	}
}

// Start will start the stream with the given buffer size.
// If the stream is already started it will be stopped then restarted with the new buffer size.
func (s *Stream) Start(bufferSize int) {
	_ = s.close()
	s.l.Lock()
	defer s.l.Unlock()
	s.stop()
	s.c = make(chan interface{}, bufferSize)
	s.started = true
	s.stopchan = make(chan struct{})
}

func (s *Stream) close() bool {
	select {
	case <-s.stopchan:
		return false
	default:
		s.stoplock.Lock()
		defer s.stoplock.Unlock()

		if s.stopchan == nil {
			return false
		}

		close(s.stopchan)

		return true
	}
}

// Stop will stop the stream and close the underlying channel
// returns false if the stream was already stopped.
func (s *Stream) Stop() bool {
	if !s.close() {
		return false
	}

	s.l.Lock()
	defer s.l.Unlock()

	return s.stop()
}

// C returns a receiver channel for the stream if the stream is stopped a stopchan channel will be returned.
func (s *Stream) C() <-chan interface{} {
	s.l.Lock()
	defer s.l.Unlock()

	if s.c != nil {
		return s.c
	}

	closedChannel := make(chan interface{})
	close(closedChannel)

	return closedChannel
}

// stop closes the underlying channel.
func (s *Stream) stop() bool {
	if !s.started {
		return false
	}
	// s.c is never nil because the stream was stareted
	close(s.c)
	s.started = false

	return true
}

// IsStopped returns true if the stream is currently stopped.
func (s *Stream) IsStopped() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return s.isStopped()
}

func (s *Stream) isStopped() bool { return !s.started }

// IsZero returns true if the stream is zero intialized.
func (s *Stream) IsZero() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return s.c == nil
}

// Send sends the given value to the stream. If the channel is unable to send then the oldest value will be discarded.
// For un buffered channels if there are no receives the function will return immediately
// There is a chance if you send on a full stream
// two values from the stream will be discarded instaead of just one.
// If the stream is stopped this function returns immediately.
func (s *Stream) Send(value interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	if s.isStopped() {
		return
	}

	select {
	case s.c <- value:
		return
	default: // do nothing
	}

	if cap(s.c) == 0 {
		return
	}

	s.loopUntilSend(value)
}

func (s *Stream) loopUntilSend(value interface{}) {
	// separate select statements for sending and receiveing so that we guarantee that the send is tried first,
	// and don't risk draining all values from the channel.
	// since there is not priority when determining which select case is run if both are true.
	const checkDelay = 10
	ticker := time.NewTicker(time.Millisecond * checkDelay)
	for {
		select {
		case s.c <- value:
			return
		default:
		}
		select {
		// read off the oldest value to make room for a send
		case <-s.c:
			// try to send again
			select {
			case s.c <- value:
				return
			default: // do nothing
			}
		case <-ticker.C:
		}
	}
}

// SendBlocking sends the given value to the stream
// if the stream is full then the send will block with the same behavior as sending on a channel that is full.
// If the stream is stopped this function returns immediately.
// If used concurrently with send then the function may not block as expected.
func (s *Stream) SendBlocking(value interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	if s.isStopped() {
		return
	}

	select {
	case s.c <- value:
	case <-s.stopchan:
		return
	}
}
