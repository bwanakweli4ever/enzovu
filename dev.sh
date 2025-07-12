#!/bin/bash
echo "🔄 Watching for changes..."

while true; do
    echo "🚀 Running: go run main.go"
    go run main.go &
    PID=$!
    
    # Wait for file changes
    inotifywait -e modify -r . --include='\.go$' 2>/dev/null || sleep 2
    
    echo "📝 File changed, restarting..."
    kill $PID 2>/dev/null
    sleep 1
done