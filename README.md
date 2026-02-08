# Redis Implementation with Go

A lightweight Redis-like in-memory key-value store built from scratch in Go using raw TCP sockets.

## Architecture

The project follows Go's standard project layout with clean separation of concerns:

```
redis-go/
├── cmd/
│   └── server/
│       └── main.go              # Entry point - wires all dependencies
├── internal/
│   ├── commands/
│   │   └── register.go          # Command definitions (PING, GET, SET)
│   ├── router/
│   │   └── router.go            # Command routing and dispatching
│   └── store/
│       └── store.go             # Thread-safe in-memory key-value store
├── pkg/
│   └── connection/
│       └── server.go            # TCP server and connection handling
├── go.mod
└── README.md
```

### Package Responsibilities

| Package | Role |
|---------|------|
| `cmd/server` | Application entry point. Creates dependencies and wires them together. |
| `internal/router` | Maps command names to handler functions. Parses incoming messages. |
| `internal/commands` | Defines Redis command handlers (PING, GET, SET). Receives store via dependency injection. |
| `internal/store` | Thread-safe `map[string]string` with `sync.RWMutex` for concurrent access. |
| `pkg/connection` | TCP listener that accepts connections and delegates messages to a handler. |

### Request Flow

```
Client (TCP) → TCPServer → Router.Handle → CommandHandler → Store → Response
```

1. Client sends a text command over TCP (e.g. `SET name redis\n`)
2. `TCPServer` reads the line and passes it to `Router.Handle`
3. `Router` parses the command name, finds the matching handler
4. Handler executes the logic (reads/writes to `Store`)
5. Response is sent back to the client

### Dependency Injection

Dependencies are wired explicitly in `main.go` without any framework:

```go
s := store.New()          // Create store
r := router.New()         // Create router
commands.Register(r, s)   // Register commands with store injected
```

## Supported Commands

| Command | Usage | Description |
|---------|-------|-------------|
| `PING` | `PING` | Returns `PONG`. Used for health checks. |
| `GET` | `GET <key>` | Returns the value for the given key, or `(nil)` if not found. |
| `SET` | `SET <key> <value>` | Stores the key-value pair. Returns `OK`. |

## Getting Started

### Prerequisites

- Go 1.21+

### Run the Server

```bash
go run ./cmd/server
```

Server starts listening on port `8080`.

### Connect with a Client

Using `nc` (netcat):

```bash
nc localhost 8080
```

Example session:

```
PING
PONG
SET name redis
OK
GET name
redis
GET nonexistent
(nil)
```

Using `telnet`:

```bash
telnet localhost 8080
```

### Build

```bash
go build -o redis-server ./cmd/server
./redis-server
```

## Concurrency

The store uses `sync.RWMutex` for thread safety:

- **Read operations** (`GET`) use `RLock` - multiple readers can access simultaneously
- **Write operations** (`SET`) use `Lock` - exclusive access, blocks all other reads and writes

Each client connection runs in its own goroutine, so multiple clients can connect and send commands concurrently.
