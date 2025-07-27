.PHONY: clean linux windows mac all help

# Default target
all: help

# Show help
help:
	@echo "Available targets:"
	@echo "  clean   - Remove all files from bin directory"
	@echo "  linux   - Build binaries for all Linux architectures"
	@echo "  windows - Build binaries for all Windows architectures"
	@echo "  mac     - Build binaries for all macOS architectures"
	@echo "  all     - Build binaries for all platforms"
	@echo "  help    - Show this help message"

# Clean build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf bin/*
	@echo "✓ Build directory cleaned"

# Build for Linux (all architectures)
linux:
	@echo "Building for Linux..."
	@./build.sh --goos linux
	@echo "✓ Linux builds complete"

# Build for Windows (all architectures)
windows:
	@echo "Building for Windows..."
	@./build.sh --goos windows
	@echo "✓ Windows builds complete"

# Build for macOS (all architectures)
mac:
	@echo "Building for macOS..."
	@./build.sh --goos darwin
	@echo "✓ macOS builds complete"

# Build for all platforms
all-platforms:
	@echo "Building for all platforms..."
	@./build.sh
	@echo "✓ All platform builds complete"
