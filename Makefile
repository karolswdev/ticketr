.PHONY: poc poc-demo build-poc help

# Build and run the Bubbletea TUI POC
poc:
	@echo "🎫 Starting Ticketr Bubbletea POC (normal mode)..."
	@go run cmd/ticketr-tui-poc/main.go

# Run POC in demo mode (themes cycle automatically)
poc-demo:
	@echo "🎫 Starting Ticketr Bubbletea POC (demo mode)..."
	@echo "   Themes will cycle every 3 seconds"
	@echo "   Press 'q' or Ctrl+C to quit"
	@echo ""
	@go run cmd/ticketr-tui-poc/main.go -demo

# Build POC binary
build-poc:
	@echo "🔨 Building Ticketr POC binary..."
	@go build -o bin/ticketr-poc cmd/ticketr-tui-poc/main.go
	@echo "✅ Binary created: bin/ticketr-poc"

# Show help for POC targets
help:
	@echo "Ticketr Bubbletea POC - Make Targets"
	@echo ""
	@echo "  make poc          - Run POC in normal mode"
	@echo "  make poc-demo     - Run POC in demo mode (auto theme cycling)"
	@echo "  make build-poc    - Build POC binary to bin/ticketr-poc"
	@echo "  make help         - Show this help message"
	@echo ""
	@echo "Keyboard shortcuts in POC:"
	@echo "  Tab       - Switch focus between panels"
	@echo "  1/2/3     - Switch themes (Default/Dark/Arctic)"
	@echo "  q, Ctrl+C - Quit"
	@echo ""
