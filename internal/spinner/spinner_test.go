package spinner

import (
	"bytes"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	sp := New("test")
	if sp == nil {
		t.Fatal("spinner should not be nil")
	}
	if sp.message != "test" {
		t.Error("message not set correctly")
	}
	if sp.running {
		t.Error("should not be running initially")
	}
}

func TestStartStop(t *testing.T) {
	var buf bytes.Buffer
	sp := New("testing").WithWriter(&buf)

	if sp.IsRunning() {
		t.Error("should not be running before Start")
	}

	sp.Start()
	if !sp.IsRunning() {
		t.Error("should be running after Start")
	}

	time.Sleep(150 * time.Millisecond)
	sp.Stop()

	if sp.IsRunning() {
		t.Error("should not be running after Stop")
	}

	if buf.Len() == 0 {
		t.Error("expected output to be produced")
	}
}

func TestUpdateMessage(t *testing.T) {
	var buf bytes.Buffer
	sp := New("first").WithWriter(&buf)

	sp.Start()
	time.Sleep(100 * time.Millisecond)
	sp.UpdateMessage("second")
	time.Sleep(100 * time.Millisecond)
	sp.Stop()

	output := buf.String()
	if output == "" {
		t.Error("expected output")
	}
}

func TestFrames(t *testing.T) {
	frames := []FrameSet{Dots, Dots2, Line, Arrow, Box, Circle}
	for i, f := range frames {
		if len(f) == 0 {
			t.Errorf("frame set %d is empty", i)
		}
	}
}
