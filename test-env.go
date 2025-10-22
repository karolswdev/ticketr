package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
)

func main() {
	fmt.Println("======================================")
	fmt.Println("Environment Variable Test")
	fmt.Println("======================================")

	fmt.Printf("TICKETR_THEME = %q\n", os.Getenv("TICKETR_THEME"))
	fmt.Printf("TICKETR_EFFECTS_AMBIENT = %q\n", os.Getenv("TICKETR_EFFECTS_AMBIENT"))
	fmt.Printf("TICKETR_EFFECTS_MOTION = %q\n", os.Getenv("TICKETR_EFFECTS_MOTION"))

	fmt.Println()
	fmt.Println("Loading theme from environment...")
	theme.LoadThemeFromEnv()

	current := theme.Current()
	fmt.Printf("Current theme: %s\n", current.Name)

	effects := theme.GetEffects()
	fmt.Printf("Motion: %v\n", effects.Motion)
	fmt.Printf("AmbientEnabled: %v\n", effects.AmbientEnabled)
	fmt.Printf("AmbientMode: %s\n", effects.AmbientMode)
	fmt.Printf("AmbientDensity: %f\n", effects.AmbientDensity)
	fmt.Printf("AmbientSpeed: %d\n", effects.AmbientSpeed)

	fmt.Println()
	fmt.Println("Parsing TICKETR_EFFECTS_AMBIENT manually...")
	ambientStr := os.Getenv("TICKETR_EFFECTS_AMBIENT")
	if ambient, err := strconv.ParseBool(ambientStr); err == nil {
		fmt.Printf("Parsed value: %v\n", ambient)
	} else {
		fmt.Printf("Parse error: %v\n", err)
	}
}
