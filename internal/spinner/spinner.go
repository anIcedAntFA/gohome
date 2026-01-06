// Package spinner provides terminal spinner animations for indicating progress.
package spinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Spinner represents a terminal spinner with customizable frames and behavior.
type Spinner struct {
	frames   FrameSet
	interval time.Duration
	message  string
	writer   io.Writer
	done     chan bool
	mu       sync.Mutex
	running  bool
}

// New creates a new Spinner with the given message using default settings.
func New(message string) *Spinner {
	return &Spinner{
		frames:   Dots,
		interval: 80 * time.Millisecond,
		message:  message,
		writer:   os.Stderr,
		done:     make(chan bool),
		running:  false,
	}
}

// WithFrames sets custom animation frames for the spinner.
func (s *Spinner) WithFrames(frames FrameSet) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = frames
	return s
}

// WithInterval sets the frame update interval for the spinner.
func (s *Spinner) WithInterval(interval time.Duration) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interval = interval
	return s
}

// WithWriter sets the output writer for the spinner.
func (s *Spinner) WithWriter(w io.Writer) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writer = w
	return s
}

// Start begins the spinner animation in a separate goroutine.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.done = make(chan bool)
	s.mu.Unlock()

	go s.animate()
}

// animate runs the animation loop in a goroutine.
func (s *Spinner) animate() {
	fmt.Fprint(s.writer, "\033[?25l")

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	frameIndex := 0
	for {
		select {
		case <-s.done:
			fmt.Fprint(s.writer, "\r\033[K\033[?25h")
			return
		case <-ticker.C:
			s.mu.Lock()
			frame := s.frames[frameIndex%len(s.frames)]
			msg := s.message
			s.mu.Unlock()

			fmt.Fprintf(s.writer, "\r%s %s", frame, msg)
			frameIndex++
		}
	}
}

// Stop stops the spinner animation and cleans up the terminal.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	s.done <- true
	close(s.done)
}

// UpdateMessage updates the spinner message while it is running.
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// IsRunning returns whether the spinner is currently running.
func (s *Spinner) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
