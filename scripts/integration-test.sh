#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TEST_PROJECT_DIR="${TEST_PROJECT_DIR:-/tmp/gogo-integration-test}"
CONTAINER_NAME="gogo-integration-postgres"

# Function to check if port is in use (works on both macOS and Linux)
port_in_use() {
    local port=$1
    if command -v lsof > /dev/null 2>&1; then
        lsof -i ":$port" > /dev/null 2>&1
    elif command -v ss > /dev/null 2>&1; then
        ss -tuln | grep -q ":$port "
    else
        # Fallback: try to connect
        (echo > /dev/tcp/localhost/$port) 2>/dev/null
    fi
}

# Check if running in CI (GitHub Actions sets CI=true)
if [ "$CI" = "true" ]; then
    echo "Running in CI environment"
    # In CI, postgres service is already running on port 5432
    PG_PORT="${POSTGRES_PORT:-5432}"
    SKIP_POSTGRES_CONTAINER=true
else
    echo "Running in local environment"
    # Find available postgres port locally
    PG_PORT=5433
    for p in $(seq 5433 5450); do
        if ! port_in_use "$p"; then
            PG_PORT=$p
            break
        fi
    done
    SKIP_POSTGRES_CONTAINER=false
fi
echo "Using PostgreSQL port: $PG_PORT"

cleanup() {
    echo "Cleaning up..."
    if [ -f /tmp/gogo-app.pid ]; then
        kill $(cat /tmp/gogo-app.pid) 2>/dev/null || true
        rm -f /tmp/gogo-app.pid
    fi
    if [ "$SKIP_POSTGRES_CONTAINER" = "false" ]; then
        docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
    fi
    rm -rf "$TEST_PROJECT_DIR"
}

trap cleanup EXIT

echo "=== Starting integration test ==="

# Clean up any previous run
cleanup 2>/dev/null || true

if [ "$SKIP_POSTGRES_CONTAINER" = "false" ]; then
    # Start postgres container (local only)
    echo "Starting PostgreSQL container..."
    docker run -d \
        --name "$CONTAINER_NAME" \
        -e POSTGRES_USER=testuser \
        -e POSTGRES_PASSWORD=testpass \
        -e POSTGRES_DB=testdb \
        -p "$PG_PORT":5432 \
        postgres:16

    # Wait for postgres to be ready
    echo "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if docker exec "$CONTAINER_NAME" pg_isready -U testuser -d testdb > /dev/null 2>&1; then
            echo "PostgreSQL is ready"
            break
        fi
        echo "Waiting... ($i/30)"
        sleep 1
    done
else
    echo "Using existing PostgreSQL service..."
fi

# Build gogo-cli
echo "Building gogo-cli..."
cd "$PROJECT_ROOT"
go build -o gogo-cli .

# Generate test project
echo "Generating test project..."
mkdir -p "$TEST_PROJECT_DIR"
./gogo-cli generate testapp \
    --email="test@example.com" \
    --repo="github.com/test/testapp" \
    --testemail="test@example.com" \
    --out="$TEST_PROJECT_DIR"

# Run migrations
echo "Running migrations..."
export POSTGRES_HOST=localhost
export POSTGRES_PORT="$PG_PORT"
export POSTGRES_DB=testdb
export POSTGRES_USER=testuser
export POSTGRES_PASSWORD=testpass

docker run --rm \
    -v "$TEST_PROJECT_DIR":/workspace \
    --network host \
    -e POSTGRES_HOST \
    -e POSTGRES_PORT \
    -e POSTGRES_DB \
    -e POSTGRES_USER \
    -e POSTGRES_PASSWORD \
    --entrypoint sql-migrate \
    ghcr.io/can3p/gogo-docker-deps:latest \
    up -env=development

# Generate models with sqlboiler
echo "Generating models with sqlboiler..."
# Substitute environment variables in sqlboiler.toml
envsubst < "$TEST_PROJECT_DIR/sqlboiler.toml" > "$TEST_PROJECT_DIR/sqlboiler.toml.tmp"
mv "$TEST_PROJECT_DIR/sqlboiler.toml.tmp" "$TEST_PROJECT_DIR/sqlboiler.toml"

docker run --rm \
    -v "$TEST_PROJECT_DIR":/workspace \
    --network host \
    ghcr.io/can3p/gogo-docker-deps:latest \
    --add-panic-variants --no-hooks --no-tests --add-enum-types \
    -c sqlboiler.toml psql

# Install frontend dependencies and build
echo "Building frontend..."
cd "$TEST_PROJECT_DIR/cmd/web"
yarn install
yarn build

# Download Go dependencies
echo "Downloading Go dependencies..."
cd "$TEST_PROJECT_DIR"
go mod download

# Find available port for the app
echo "Finding available port..."
PORT=8080
for p in $(seq 8080 8180); do
    if ! port_in_use "$p"; then
        PORT=$p
        break
    fi
done
echo "Using port: $PORT"

# Start application
echo "Starting application..."
cd "$TEST_PROJECT_DIR/cmd/web"
export DATABASE_URL="postgres://testuser:testpass@localhost:$PG_PORT/testdb?sslmode=disable"
export SESSION_SALT="testsalt12345678"
export SITE_ROOT="http://localhost:$PORT"
export GIN_MODE=release
export PORT

go run . &
echo $! > /tmp/gogo-app.pid

# Wait for app to start
echo "Waiting for application to start..."
for i in {1..30}; do
    if curl -s "http://localhost:$PORT" > /dev/null 2>&1; then
        echo "Application started successfully on port $PORT"
        break
    fi
    echo "Waiting... ($i/30)"
    sleep 1
done

# Test application health
echo "Testing application health..."
response=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT")
if [ "$response" != "200" ]; then
    echo "❌ Health check failed with status: $response"
    exit 1
fi
echo "✅ Health check passed with status: $response"

# Test login page
echo "Testing login page..."
response=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT/login")
if [ "$response" != "200" ]; then
    echo "❌ Login page check failed with status: $response"
    exit 1
fi
echo "✅ Login page check passed with status: $response"

echo ""
echo "=== Integration test PASSED ==="
