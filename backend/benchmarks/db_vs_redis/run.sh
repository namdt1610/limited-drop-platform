#!/bin/bash

echo "🚀 Setting up Benchmark Lab..."

# Check if we are in Codespaces (Redis is already running)
if [ "$CODESPACES" == "true" ]; then
    echo "☁️  Running in Cloud Codespaces..."
else
    echo "💻 Running locally..."
    # Check if redis container is running
    if [ ! "$(docker ps -q -f name=bench-redis)" ]; then
        if [ "$(docker ps -aq -f name=bench-redis)" ]; then
            echo "🔄 Restarting existing redis container..."
            docker start bench-redis
        else
            echo "📦 Starting new redis container..."
            docker run --rm -d -p 6379:6379 --name bench-redis redis:alpine
        fi
        
        # Postgres setup
        if [ ! "$(docker ps -q -f name=bench-pg)" ]; then
            if [ "$(docker ps -aq -f name=bench-pg)" ]; then
                 echo "🔄 Restarting existing postgres container..."
                 docker start bench-pg
            else
                 echo "🐘 Starting new postgres container..."
                 docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name bench-pg postgres:alpine
            fi
            # Wait for postgres to be ready
            echo "⏳ Waiting for Postgres..."
            sleep 5
        fi

        # Wait for services to be ready
        sleep 2
    fi
fi

echo "🏁 Running Benchmark..."
go run benchmarks/db_vs_redis/main.go

# Cleanup only if local (keep codespaces persistent)
if [ "$CODESPACES" != "true" ]; then
    echo "🧹 Cleaning up..."
    docker stop bench-redis
    docker stop bench-pg
    rm benchmark.db
fi
