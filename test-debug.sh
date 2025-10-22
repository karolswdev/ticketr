#!/bin/bash

export TICKETR_THEME=dark
export TICKETR_EFFECTS_AMBIENT=true

echo "Starting TUI with debug logging..."
echo "Press 'W' to toggle workspace panel"
echo "Press 'q' to quit"
echo ""

timeout 10s ./ticketr tui 2>&1 | tee cosmic-debug.log || true

echo ""
echo "========================================"
echo "Debug output:"
echo "========================================"
grep "\[DEBUG" cosmic-debug.log | head -30
