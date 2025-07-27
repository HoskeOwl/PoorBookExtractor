.PHONY: clean linux windows mac all help version

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
	@echo "  version - Show current version information"
	@echo "  help    - Show this help message"

# Show version information
version:
	@echo "Current version information:"
	@echo "Git tag: $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'unknown')"
	@echo "Git commit: $(shell git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo 'unknown')"
	@echo "Build time: $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')"

# Clean build directory
clean:
	@echo "Cleaning build directory..."
	@rm -rf bin/*
	@echo "✓ Build directory cleaned"

# Build for Linux (all architectures)
linux:
	@echo "Building for Linux..."
	@chmod +x build.sh
	@./build.sh --goos linux
	@echo "✓ Linux builds complete"

# Build for Windows (all architectures)
windows:
	@echo "Building for Windows..."
	@chmod +x build.sh
	@./build.sh --goos windows
	@echo "✓ Windows builds complete"

# Build for macOS (all architectures)
mac:
	@echo "Building for macOS..."
	@chmod +x build.sh
	@./build.sh --goos darwin
	@echo "✓ macOS builds complete"

# Build for all platforms
all-platforms:
	@echo "Building for all platforms..."
	@chmod +x build.sh
	@./build.sh
	@echo "✓ All platform builds complete"
