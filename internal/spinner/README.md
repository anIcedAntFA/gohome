# Spinner Package

Internal package providing terminal spinner animations for indicating progress.

## Features

- 🎨 **15+ Built-in Frame Sets** - From classic spinners to fun animations
- ⚡ **Zero Dependencies** - Pure Go implementation
- 🔒 **Thread-Safe** - Goroutine-safe with mutex protection
- 🎛️ **Fully Customizable** - Frames, interval, output writer
- 💻 **Cross-Platform** - Works on Linux, macOS, Windows

## Quick Start

```go
import "github.com/anIcedAntFA/gohome/internal/spinner"

// Basic usage
sp := spinner.New("Loading...")
sp.Start()
// ... do work ...
sp.Stop()
```

## Available Frame Sets

### Classic Spinners

- **Dots** (default) - `⠋ ⠙ ⠹ ⠸` - Braille pattern spinner
- **Dots2** - `⣾ ⣽ ⣻ ⢿` - Alternative Braille pattern
- **Dots3** - `⣷ ⣯ ⣟ ⡿` - Triple dots variation
- **Line** - `- \ | /` - Classic rotating line
- **Circle** - `◴ ◷ ◶ ◵` - Circular animation
- **Box** - `◰ ◳ ◲ ◱` - Box drawing characters

### Directional & Progress

- **Arrow** - `← ↖ ↑ ↗ →` - Directional arrows
- **ProgressBar** - `▱▱▱ ▰▱▱ ▰▰▱` - Loading bar
- **BouncingBar** - `[= ]  [== ]` - Bouncing progress bar
- **Wave** - `◜ ◝ ◞ ◟` - Wave/pulse animation
- **Bounce** - `⠁ ⠂ ⠄ ⡀` - Bouncing ball

### Fun & Creative

- **Pacman** - `ᗧ··· ᗣ·· ᗣᗣ·` - Classic Pacman chomping
- **Shark** - `▐|\\___` - Swimming shark animation
- **Earth** - `🌍 🌎 🌏` - Rotating earth
- **Moon** - `🌑 🌒 🌓 🌔` - Moon phases
- **Clock** - `🕐 🕑 🕒 🕓` - Clock rotation
- **Hearts** - `💛 💙 💜 💚 ❤️` - Beating hearts

## Usage Examples

### 1. Basic with Default Settings

```go
sp := spinner.New("Processing...")
sp.Start()
time.Sleep(2 * time.Second)
sp.Stop()
fmt.Println("✓ Done!")
```

### 2. Custom Frames (Pacman)

```go
sp := spinner.New("Chomping files...").WithFrames(spinner.Pacman)
sp.Start()
time.Sleep(2 * time.Second)
sp.Stop()
```

### 3. Custom Speed (Faster)

```go
sp := spinner.New("Quick task...").
    WithFrames(spinner.Wave).
    WithInterval(50 * time.Millisecond)  // Default is 80ms

sp.Start()
time.Sleep(2 * time.Second)
sp.Stop()
```

### 4. Update Message While Running

```go
sp := spinner.New("Step 1...").WithFrames(spinner.Clock)
sp.Start()

time.Sleep(1 * time.Second)
sp.UpdateMessage("Step 2...")

time.Sleep(1 * time.Second)
sp.UpdateMessage("Step 3...")

time.Sleep(1 * time.Second)
sp.Stop()
```

### 5. Custom Output Writer

```go
var buf bytes.Buffer

sp := spinner.New("Writing to buffer...").
    WithWriter(&buf).
    WithFrames(spinner.Hearts)

sp.Start()
time.Sleep(1 * time.Second)
sp.Stop()

fmt.Println(buf.String())  // See captured output
```

### 6. Multiple Spinners in Sequence

```go
// Scanning phase
sp1 := spinner.New("🔍 Scanning...").WithFrames(spinner.Dots2)
sp1.Start()
time.Sleep(1 * time.Second)
sp1.Stop()
fmt.Println("✓ Found 42 items")

// Processing phase
sp2 := spinner.New("⚙️  Processing...").WithFrames(spinner.BouncingBar)
sp2.Start()
time.Sleep(2 * time.Second)
sp2.Stop()
fmt.Println("✓ Processed")
```

### 7. Fully Customized

```go
sp := spinner.New("Processing items...").
    WithFrames(spinner.Pacman).
    WithInterval(60 * time.Millisecond).
    WithWriter(os.Stdout)  // Or any io.Writer

sp.Start()

for i := 1; i <= 5; i++ {
    time.Sleep(500 * time.Millisecond)
    sp.UpdateMessage(fmt.Sprintf("Item %d/5...", i))
}

sp.Stop()
fmt.Println("✓ Complete!")
```

## API Reference

### Constructor

```go
func New(message string) *Spinner
```

Creates a new spinner with default settings (Dots frames, 80ms interval, stderr output).

### Configuration Methods (Chainable)

```go
func (s *Spinner) WithFrames(frames FrameSet) *Spinner
```

Set custom animation frames.

```go
func (s *Spinner) WithInterval(interval time.Duration) *Spinner
```

Set frame update interval (default: 80ms).

```go
func (s *Spinner) WithWriter(w io.Writer) *Spinner
```

Set output writer (default: os.Stderr).

### Control Methods

```go
func (s *Spinner) Start()
```

Start the spinner animation (non-blocking, runs in goroutine).

```go
func (s *Spinner) Stop()
```

Stop the spinner and cleanup (hides cursor, clears line).

```go
func (s *Spinner) UpdateMessage(message string)
```

Update the message while spinner is running (thread-safe).

```go
func (s *Spinner) IsRunning() bool
```

Check if spinner is currently running.

## Tips

### Choosing the Right Spinner

- **Fast operations** (<1s): Use simple spinners like `Dots`, `Line`, or `Wave`
- **Medium operations** (1-5s): Use engaging spinners like `Pacman`, `Earth`, or `ProgressBar`
- **Long operations** (>5s): Use `BouncingBar`, `Shark`, or update messages periodically
- **Serious/Professional**: Stick to `Dots`, `Dots2`, `Line`, `Circle`
- **Fun/Casual**: Try `Pacman`, `Hearts`, `Earth`, `Moon`

### Performance

- Default interval (80ms) is optimal for most cases
- Faster intervals (50-60ms) for quick, snappy feel
- Slower intervals (100-150ms) for calmer, less distracting animation
- Shark spinner has 25 frames, best with 60-80ms interval

### Best Practices

```go
// Always defer Stop() to ensure cleanup
sp := spinner.New("Working...")
sp.Start()
defer sp.Stop()

// Your code here...
```

```go
// Check IsRunning() before stopping
if sp.IsRunning() {
    sp.Stop()
}
```

```go
// Use meaningful messages
sp.UpdateMessage("Downloading dependencies...")  // ✓ Good
sp.UpdateMessage("Working...")                   // ✗ Too vague
```

## Thread Safety

All methods are thread-safe. You can safely:

- Call `UpdateMessage()` from multiple goroutines
- Call `IsRunning()` concurrently
- Start/Stop from different goroutines

## Notes

- Spinner writes to `os.Stderr` by default (doesn't interfere with stdout)
- ANSI escape codes are used for cursor control
- Call `Stop()` before printing final output to ensure clean terminal state
- Spinner clears its line on stop, so print your completion message after

## License

Internal package - part of gohome project.
