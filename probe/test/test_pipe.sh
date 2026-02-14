#!/bin/bash
set -e

PROBE="../probe"
PIPE="/tmp/shinobi.pipe"
TIMEOUT=5

echo "🧪 Testing named pipe integration..."

# Clean up any previous test artifacts
rm -f "$PIPE"
pkill -f "$PROBE" 2>/dev/null || true

# Start probe in background
"$PROBE" &
PROBE_PID=$!

# Give it a moment to start
sleep 1

# Read one line from the pipe with timeout (macOS compatible)
echo "Reading from $PIPE..."
LINE=$(perl -e 'alarm shift @ARGV; exec @ARGV' "$TIMEOUT" head -n 1 "$PIPE" 2>/dev/null || true)

if [ -z "$LINE" ]; then
    echo "❌ Failed to read from pipe"
    kill "$PROBE_PID" 2>/dev/null || true
    exit 1
fi

echo "Received: $LINE"

# Validate JSON structure
if ! echo "$LINE" | grep -q '"cpu_percent"'; then
    echo "❌ Missing 'cpu_percent' key"
    kill "$PROBE_PID" 2>/dev/null || true
    exit 1
fi

if ! echo "$LINE" | grep -q '"timestamp"'; then
    echo "❌ Missing 'timestamp' key"
    kill "$PROBE_PID" 2>/dev/null || true
    exit 1
fi

# Extract cpu_percent and validate range
CPU=$(echo "$LINE" | grep -o '"cpu_percent":[0-9.]*' | cut -d: -f2)
if ! echo "$CPU" | grep -qE '^[0-9]+(\.[0-9]+)?$'; then
    echo "❌ Invalid cpu_percent value: $CPU"
    kill "$PROBE_PID" 2>/dev/null || true
    exit 1
fi

# Check if value is in valid range (0-100)
if (( $(echo "$CPU < 0" | bc -l) )) || (( $(echo "$CPU > 100" | bc -l) )); then
    echo "❌ cpu_percent out of range: $CPU"
    kill "$PROBE_PID" 2>/dev/null || true
    exit 1
fi

echo "✓ Valid JSON with cpu_percent=$CPU%"

# Clean up
kill "$PROBE_PID" 2>/dev/null || true
rm -f "$PIPE"

echo "✓ Integration test passed!"
