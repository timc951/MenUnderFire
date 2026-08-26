#!/bin/bash

env=$1

# Determine which .env file to use
env_file=".env"
if [ ! -f "$env_file" ]; then
    echo "ERROR: $env_file file not found"
    exit 1
fi

# Check for required environment variables
required_file=".env.required"
if [ -f "$required_file" ]; then
    missing=()
    while IFS= read -r var || [ -n "$var" ]; do
        # Skip empty lines and comments
        var=$(echo "$var" | xargs)
        [[ -z "$var" || "$var" == \#* ]] && continue

        # Check if the variable is set (non-empty) in the .env file
        if ! grep -q "^${var}=.\+" "$env_file"; then
            missing+=("$var")
        fi
    done < "$required_file"

    if [ ${#missing[@]} -gt 0 ]; then
        echo "ERROR: Missing required environment variables in $env_file:"
        for var in "${missing[@]}"; do
            echo "  - $var"
        done
        exit 1
    fi
fi

if [ -n "$env" ]; then
    docker compose -f docker-compose.$env.yml build --no-cache && docker compose -f docker-compose.$env.yml up -d
else
    docker compose build --no-cache && docker compose up -d
fi
