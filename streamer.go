// Package streamer is used to stream data utilizing channels
package streamer

import "sync"

type Streamer struct {
	c       chan interface{}
	l       sync.Mutex
	started bool
}

// New reuturns a new streamer object with the specified buffer size.
func New(bufferSize int) *Streamer {
	return &Streamer{
		c:       make(chan interface{}, bufferSize),
		started: true,
		l:       sync.Mutex{},
	}
}

// Start will start the streamer with the given buffer size.
// If the streamer is already started it will be stopped then restarted with the new buffer size.
func (s *Streamer) Start(bufferSize int) {
	s.l.Lock()
	defer s.l.Unlock()
	_ = s.stop()
	s.c = make(chan interface{}, bufferSize)
	s.started = true
}

// Stop will stop the streamer and close the underlying channel
// returns false if the stream was already stoped.
func (s *Streamer) Stop() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return s.stop()
}

// C returns a receiver channel for the streamer if the streamer is stopped a closed channel will be returned.
func (s *Streamer) C() <-chan interface{} {
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
func (s *Streamer) stop() bool {
	if s.started {
		return false
	}

	close(s.c)
	s.started = false

	return true
}

// IsStopped returns true if the streamer is currently stopped
func (s *Streamer) IsStopped() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return !s.started
}

// IsZero returns true if the streamer is zero intialized
func (s *Streamer) IsZero() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return s.c == nil
}
