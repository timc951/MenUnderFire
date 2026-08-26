#!/bin/bash

# Script to check if the golang docker container image is the latest version
# Compares Dockerfile's golang version against Docker Hub's latest

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default Dockerfile path (look in parent directory)
DOCKERFILE="${1:-../Dockerfile}"

# Try alternate locations if not found
if [ ! -f "$DOCKERFILE" ]; then
    if [ -f "Dockerfile" ]; then
        DOCKERFILE="Dockerfile"
    elif [ -f "../../Dockerfile" ]; then
        DOCKERFILE="../../Dockerfile"
    else
        echo -e "${RED}Error: Dockerfile not found${NC}"
        echo "Usage: $0 [path/to/Dockerfile]"
        exit 1
    fi
fi

# Extract golang version from Dockerfile (matches patterns like golang:1.24.13-alpine3.23)
CURRENT_IMAGE=$(grep -oP 'FROM golang:\K[^\s]+' "$DOCKERFILE" | head -1)

if [ -z "$CURRENT_IMAGE" ]; then
    echo -e "${RED}Error: Could not find golang image in Dockerfile${NC}"
    exit 1
fi

# Extract just the Go version number (e.g., 1.24.13 from 1.24.13-alpine3.23)
CURRENT_VERSION=$(echo "$CURRENT_IMAGE" | grep -oP '^[0-9]+\.[0-9]+(\.[0-9]+)?')

# Extract the variant (e.g., alpine3.23, bookworm, etc.)
CURRENT_VARIANT=$(echo "$CURRENT_IMAGE" | grep -oP '(?<=-).+$' || echo "")

echo "Dockerfile: $DOCKERFILE"
echo "Current golang image: golang:$CURRENT_IMAGE"
echo "Current Go version: $CURRENT_VERSION"
echo "Current variant: ${CURRENT_VARIANT:-<none>}"
echo ""

# Query Docker Hub API for golang image tags
echo "Fetching latest golang versions from Docker Hub..."

# Get the list of tags from Docker Hub (official golang image)
TAGS_JSON=$(curl -s "https://hub.docker.com/v2/repositories/library/golang/tags?page_size=100&ordering=last_updated")

if [ -z "$TAGS_JSON" ]; then
    echo -e "${RED}Error: Could not fetch tags from Docker Hub${NC}"
    exit 1
fi

LATEST_WITH_VARIANT=""

# Find latest version with same variant
if [ -n "$CURRENT_VARIANT" ]; then
    LATEST_WITH_VARIANT=$(echo "$TAGS_JSON" | \
        grep -oP '"name":\s*"\K[0-9]+\.[0-9]+(\.[0-9]+)?-'"$CURRENT_VARIANT"'(?=")' | \
        sort -V | \
        tail -1)
fi

# Get the latest stable version (X.Y.Z format, no rc/beta/alpha)
LATEST_STABLE=$(echo "$TAGS_JSON" | \
    grep -oP '"name":\s*"\K[0-9]+\.[0-9]+\.[0-9]+(?=")' | \
    grep -v -E '(rc|beta|alpha)' | \
    sort -V | \
    tail -1)

if [ -z "$LATEST_STABLE" ]; then
    echo -e "${YELLOW}Warning: Could not determine latest stable version${NC}"
    exit 0
fi

# Use variant version for comparison if available
if [ -n "$LATEST_WITH_VARIANT" ]; then
    LATEST_VERSION=$(echo "$LATEST_WITH_VARIANT" | grep -oP '^[0-9]+\.[0-9]+(\.[0-9]+)?')
    LATEST_IMAGE="$LATEST_WITH_VARIANT"
else
    LATEST_VERSION="$LATEST_STABLE"
    LATEST_IMAGE="$LATEST_STABLE"
fi

echo "Latest stable Go version: $LATEST_STABLE"
if [ -n "$CURRENT_VARIANT" ]; then
    echo "Latest with variant ($CURRENT_VARIANT): ${LATEST_WITH_VARIANT:-not available}"
fi

# Compare versions
compare_versions() {
    local v1=$1
    local v2=$2

    # Normalize versions to X.Y.Z format
    v1=$(echo "$v1" | awk -F. '{printf "%d.%d.%d", $1, $2, ($3 ? $3 : 0)}')
    v2=$(echo "$v2" | awk -F. '{printf "%d.%d.%d", $1, $2, ($3 ? $3 : 0)}')

    if [ "$v1" = "$v2" ]; then
        return 0  # Equal
    fi

    # Compare using sort -V
    local higher=$(echo -e "$v1\n$v2" | sort -V | tail -1)
    if [ "$higher" = "$v1" ]; then
        return 1  # v1 is higher
    else
        return 2  # v2 is higher
    fi
}

echo ""
echo "============================================"

compare_versions "$CURRENT_VERSION" "$LATEST_VERSION"
result=$?

case $result in
    0)
        echo -e "${GREEN}✓ You are using the latest Go version ($CURRENT_VERSION)${NC}"
        ;;
    1)
        echo -e "${YELLOW}! Your Go version ($CURRENT_VERSION) is newer than Docker Hub's latest ($LATEST_VERSION)${NC}"
        echo "  This might be a pre-release or Docker Hub hasn't updated yet."
        ;;
    2)
        echo -e "${RED}✗ Update available!${NC}"
        echo ""
        echo "  Current: golang:$CURRENT_IMAGE (Go $CURRENT_VERSION)"
        echo "  Latest:  golang:$LATEST_IMAGE (Go $LATEST_VERSION)"
        echo ""
        echo "To update your Dockerfile, change:"
        echo "  FROM golang:$CURRENT_IMAGE"
        echo "to:"
        echo "  FROM golang:$LATEST_IMAGE"
        ;;
esac

echo "============================================"
