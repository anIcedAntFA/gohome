package spinner

// FrameSet represents a collection of animation frames
type FrameSet []string

var (
	// Dots - Braille pattern spinner (default)
	Dots = FrameSet{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// Dots2 - Alternative Braille pattern
	Dots2 = FrameSet{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

	// Line - Classic rotating line
	Line = FrameSet{"-", "\\", "|", "/"}

	// Arrow - Directional arrows
	Arrow = FrameSet{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	// Box - Box drawing characters
	Box = FrameSet{"◰", "◳", "◲", "◱"}

	// Circle - Circular animation
	Circle = FrameSet{"◴", "◷", "◶", "◵"}

	// ProgressBar - Loading bar animation
	ProgressBar = FrameSet{
		"▱▱▱▱▱▱▱",
		"▰▱▱▱▱▱▱",
		"▰▰▱▱▱▱▱",
		"▰▰▰▱▱▱▱",
		"▰▰▰▰▱▱▱",
		"▰▰▰▰▰▱▱",
		"▰▰▰▰▰▰▱",
		"▰▰▰▰▰▰▰",
	}

	// Pacman - Inspired by Arch Linux pacman package manager
	Pacman = FrameSet{
		"ᗧ······",
		"ᗣ·····",
		" ᗧ····",
		" ᗣ···",
		"  ᗧ··",
		"  ᗣ·",
		"   ᗧ",
		"   ᗣ",
	}

	// PacmanGhost - Pacman being chased by ghost
	PacmanGhost = FrameSet{
		"ᗧ····  ᗣ",
		" ᗧ···  ᗣ",
		" ᗣ···ᗣ ",
		"  ᗧ··ᗣ ",
		"  ᗣ·ᗣ  ",
		"   ᗧᗣ  ",
		"   ᗣ   ",
		"  ᗧ    ",
		" ᗧ     ",
		"ᗧ      ",
	}

	// Wave - Wave/pulse animation
	Wave = FrameSet{"◜", "◝", "◞", "◟"}

	// Bounce - Bouncing ball
	Bounce = FrameSet{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}

	// Shark - Swimming shark
	Shark = FrameSet{
		"▐|\\____________▌",
		"▐_|\\___________▌",
		"▐__|\\__________▌",
		"▐___|\\_________▌",
		"▐____|\\________▌",
		"▐_____|\\_______▌",
		"▐______|\\______▌",
		"▐_______|\\_____▌",
		"▐________|\\____▌",
		"▐_________|\\___▌",
		"▐__________|\\__▌",
		"▐___________|\\_▌",
		"▐____________|\\▌",
		"▐____________/|▌",
		"▐___________/|_▌",
		"▐__________/|__▌",
		"▐_________/|___▌",
		"▐________/|____▌",
		"▐_______/|_____▌",
		"▐______/|______▌",
		"▐_____/|_______▌",
		"▐____/|________▌",
		"▐___/|_________▌",
		"▐__/|__________▌",
		"▐_/|___________▌",
	}

	// Earth - Rotating earth
	Earth = FrameSet{"🌍", "🌎", "🌏"}

	// Moon - Moon phases
	Moon = FrameSet{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}

	// Clock - Clock rotation
	Clock = FrameSet{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"}

	// Hearts - Beating hearts
	Hearts = FrameSet{"💛", "💙", "💜", "💚", "❤️"}

	// Dots3 - Triple dots
	Dots3 = FrameSet{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}

	// BouncingBar - Bouncing bar animation
	BouncingBar = FrameSet{
		"[    ]",
		"[=   ]",
		"[==  ]",
		"[=== ]",
		"[ ===]",
		"[  ==]",
		"[   =]",
		"[    ]",
		"[   =]",
		"[  ==]",
		"[ ===]",
		"[====]",
		"[=== ]",
		"[==  ]",
		"[=   ]",
	}
)
