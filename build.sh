#!/bin/bash

OUT_DIRECTORY=bin

# Parse command line arguments
BUILD_GOOS=""
BUILD_GOARCH=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --goos)
            BUILD_GOOS="$2"
            shift 2
            ;;
        --goarch)
            BUILD_GOARCH="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [--goos GOOS] [--goarch GOARCH]"
            echo "  --goos GOOS     Build only for specific OS (e.g., linux, darwin, windows)"
            echo "  --goarch GOARCH Build only for specific architecture (e.g., amd64, arm64, 386)"
            echo "  -h, --help      Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Build for all platforms"
            echo "  $0 --goos linux       # Build for all Linux architectures"
            echo "  $0 --goarch amd64     # Build for all amd64 platforms"
            echo "  $0 --goos linux --goarch amd64  # Build only for Linux amd64"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Get version information from git
get_git_tag() {
    git describe --tags --abbrev=0 2>/dev/null || echo "unknown"
}

get_git_commit() {
    git rev-parse HEAD 2>/dev/null || echo "unknown"
}

get_build_time() {
    date -u '+%Y-%m-%d_%H:%M:%S_UTC'
}

VERSION=$(get_git_tag)
COMMIT_SHA=$(get_git_commit)
BUILD_TIME=$(get_build_time)

echo "Build information:"
echo "  Version: $VERSION"
echo "  Commit: ${COMMIT_SHA:0:8}"
echo "  Build Time: $BUILD_TIME"
echo ""


# Build matrix: GOOS GOARCH BINARY_NAME
declare -a builds=(
    "darwin amd64 poorbookextractor-darwin-amd64"
    "darwin arm64 poorbookextractor-darwin-arm64"
    "linux amd64 poorbookextractor-linux-amd64"
    "linux arm64 poorbookextractor-linux-arm64"
    "linux 386 poorbookextractor-linux-386"
    "windows amd64 poorbookextractor-windows-amd64.exe"
    "windows arm64 poorbookextractor-windows-arm64.exe"
    "windows 386 poorbookextractor-windows-386.exe"
)

# Filter builds based on command line arguments
filtered_builds=()
for build in "${builds[@]}"; do
    read -r GOOS GOARCH BINARY_NAME <<< "$build"
    
    # If specific GOOS is requested, filter by it
    if [[ -n "$BUILD_GOOS" && "$GOOS" != "$BUILD_GOOS" ]]; then
        continue
    fi
    
    # If specific GOARCH is requested, filter by it
    if [[ -n "$BUILD_GOARCH" && "$GOARCH" != "$BUILD_GOARCH" ]]; then
        continue
    fi
    
    filtered_builds+=("$build")
done

# Check if any builds match the filter
if [[ ${#filtered_builds[@]} -eq 0 ]]; then
    echo "No builds match the specified criteria: GOOS=$BUILD_GOOS, GOARCH=$BUILD_GOARCH"
    echo "Available combinations:"
    for build in "${builds[@]}"; do
        read -r GOOS GOARCH BINARY_NAME <<< "$build"
        echo "  $GOOS/$GOARCH"
    done
    exit 1
fi

echo "Building ${#filtered_builds[@]} binary(ies):"
if [[ -n "$BUILD_GOOS" ]]; then
    echo "  GOOS: $BUILD_GOOS"
fi
if [[ -n "$BUILD_GOARCH" ]]; then
    echo "  GOARCH: $BUILD_GOARCH"
fi
echo ""

# Loop through filtered build combinations
for build in "${filtered_builds[@]}"; do
    # Split the build string into variables
    read -r GOOS GOARCH BINARY_NAME <<< "$build"
    
    echo "Building for $GOOS/$GOARCH -> $BINARY_NAME"

    # Clean output directory
    rm -f "$OUT_DIRECTORY/$BINARY_NAME"
    
    # Build the binary with version information
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X github.com/HoskeOwl/PoorBookExtractor/internal/version.Version=$VERSION -X github.com/HoskeOwl/PoorBookExtractor/internal/version.CommitSHA=$COMMIT_SHA -X github.com/HoskeOwl/PoorBookExtractor/internal/version.BuildTime=\"$BUILD_TIME\"" \
        -o "$OUT_DIRECTORY/$BINARY_NAME" main.go
    
    # Check if build was successful
    if [ $? -eq 0 ]; then
        echo "✓ Successfully built $BINARY_NAME"
    else
        echo "✗ Failed to build $BINARY_NAME"
    fi
    echo ""
done

echo "Build complete! Binaries are in $OUT_DIRECTORY/"
ls -la $OUT_DIRECTORY/