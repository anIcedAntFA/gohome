// Package spinner provides examples for using the spinner package.
package spinner

import (
	"bytes"
	"fmt"
	"time"
)

// ExampleBasic demonstrates basic usage with default settings.
func ExampleBasic() {
	sp := New("Loading data...")
	sp.Start()

	// Do some work
	time.Sleep(2 * time.Second)

	sp.Stop()
	fmt.Println("✓ Done!")
}

// ExampleCustomFrames demonstrates using custom frames like Pacman.
func ExampleCustomFrames() {
	sp := New("Chomping through files...").WithFrames(Pacman)
	sp.Start()

	// Do some work
	time.Sleep(2 * time.Second)

	sp.Stop()
	fmt.Println("✓ All files processed!")
}

// ExampleCustomInterval demonstrates using a custom interval for faster animation.
func ExampleCustomInterval() {
	sp := New("Processing quickly...").
		WithFrames(Wave).
		WithInterval(50 * time.Millisecond) // Faster than default 80ms

	sp.Start()
	time.Sleep(2 * time.Second)
	sp.Stop()
	fmt.Println("✓ Complete!")
}

// ExampleCustomWriter demonstrates using a custom writer like a buffer.
func ExampleCustomWriter() {
	var buf bytes.Buffer

	sp := New("Writing to buffer...").
		WithWriter(&buf).
		WithFrames(Hearts)

	sp.Start()
	time.Sleep(1 * time.Second)
	sp.Stop()

	// Print what was written to buffer
	fmt.Println("\nBuffer contents:")
	fmt.Println(buf.String())
}

// ExampleUpdateMessage demonstrates updating the spinner message while it is running.
func ExampleUpdateMessage() {
	sp := New("Step 1: Initializing...").WithFrames(Clock)
	sp.Start()

	time.Sleep(1 * time.Second)
	sp.UpdateMessage("Step 2: Processing...")

	time.Sleep(1 * time.Second)
	sp.UpdateMessage("Step 3: Finalizing...")

	time.Sleep(1 * time.Second)
	sp.Stop()
	fmt.Println("✓ All steps completed!")
}

// ExampleMultipleSpinners demonstrates using multiple spinners in sequence.
func ExampleMultipleSpinners() {
	// First spinner - scanning
	sp1 := New("🔍 Scanning directories...").WithFrames(Dots2)
	sp1.Start()
	time.Sleep(1 * time.Second)
	sp1.Stop()
	fmt.Println("✓ Found 42 files")

	// Second spinner - processing
	sp2 := New("⚙️  Processing files...").WithFrames(BouncingBar)
	sp2.Start()
	time.Sleep(2 * time.Second)
	sp2.Stop()
	fmt.Println("✓ All files processed")

	// Third spinner - uploading
	sp3 := New("☁️  Uploading results...").WithFrames(ProgressBar)
	sp3.Start()
	time.Sleep(1500 * time.Millisecond)
	sp3.Stop()
	fmt.Println("✓ Upload complete")
}

// ExampleShowcaseFrames demonstrates different available frame sets.
func ExampleShowcaseFrames() {
	framesets := []struct {
		name   string
		frames FrameSet
	}{
		{"Dots", Dots},
		{"Pacman", Pacman},
		{"Wave", Wave},
		{"Earth", Earth},
		{"Moon", Moon},
		{"Hearts", Hearts},
		{"Shark", Shark},
		{"BouncingBar", BouncingBar},
	}

	for _, fs := range framesets {
		sp := New(fmt.Sprintf("Testing %s animation...", fs.name)).
			WithFrames(fs.frames)

		sp.Start()
		time.Sleep(1500 * time.Millisecond)
		sp.Stop()
		fmt.Printf("✓ %s complete\n", fs.name)
	}
}

// ExampleConditional demonstrates conditional spinner usage based on verbose mode.
func ExampleConditional(verbose bool) {
	var sp *Spinner

	if verbose {
		sp = New("🔄 Syncing data...").WithFrames(Earth)
		sp.Start()
	}

	// Do actual work
	time.Sleep(2 * time.Second)

	if verbose && sp != nil {
		sp.Stop()
		fmt.Println("✓ Sync completed")
	}
}

// ExampleFullyCustomized demonstrates chaining all customization options together.
func ExampleFullyCustomized() {
	sp := New("Custom spinner demonstration").
		WithFrames(Pacman).                  // Custom frames
		WithInterval(60 * time.Millisecond). // Custom speed
		WithWriter(bytes.NewBuffer(nil))     // Custom output (can be os.Stdout)

	sp.Start()

	// Simulate work
	for i := 1; i <= 5; i++ {
		time.Sleep(500 * time.Millisecond)
		sp.UpdateMessage(fmt.Sprintf("Processing item %d/5...", i))
	}

	sp.Stop()
	fmt.Println("✓ Processing complete!")
}
