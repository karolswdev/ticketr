#!/bin/bash

# Test cosmic background with debug output

export TICKETR_THEME=dark
export TICKETR_EFFECTS_AMBIENT=true

echo "========================================"
echo "Testing Cosmic Background"
echo "========================================"
echo "Theme: $TICKETR_THEME"
echo "Ambient Effects: $TICKETR_EFFECTS_AMBIENT"
echo "========================================"
echo ""
echo "Instructions:"
echo "1. Press W to open workspace panel"
echo "2. Look for stars in the background"
echo "3. Check debug output below"
echo "========================================"
echo ""

# Run with stderr to file
./ticketr tui 2>cosmic-debug.log

echo ""
echo "========================================"
echo "Debug output saved to cosmic-debug.log"
echo "========================================"
