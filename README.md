# Redis Implementation with Go

A lightweight Redis-like in-memory key-value store built from scratch in Go with TLS-encrypted TCP connections.

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
| `internal/commands` | Defines Redis command handlers (PING, GET, SET, DELETE, SET_WITH_TTL). Receives store via dependency injection. |
| `internal/store` | Thread-safe key-value store with `sync.RWMutex`, TTL support, and background cleanup. |
| `pkg/connection` | TCP/TLS listener that accepts connections and delegates messages to a handler. |

### Request Flow

```
Client (TLS/TCP) → TCPServer → Router.Handle → CommandHandler → Store → Response
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
| `GET` | `GET <key>` | Returns the value for the given key, or `(nil)` if not found. Expired keys return `(nil)`. |
| `SET` | `SET <key> <value>` | Stores the key-value pair without expiration. Returns `OK`. |
| `SET_WITH_TTL` | `SET_WITH_TTL <key> <value> <seconds>` | Stores the key-value pair with a TTL in seconds. Returns `OK (TTL: <seconds>)`. |
| `DELETE` | `DELETE <key>` | Deletes the given key. Returns `true` if key existed, `false` otherwise. |

## Getting Started

### Prerequisites

- Go 1.21+
- OpenSSL (for generating TLS certificates)

### Generate TLS Certificates

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
```

This creates `cert.pem` and `key.pem` in the project root.

### Run the Server

```bash
go run ./cmd/server
```

Server starts listening on port `8080` with TLS. To run without TLS, remove `CertFile` and `KeyFile` from `main.go`.

### Connect with a Client

Using `openssl s_client` (TLS):

```bash
openssl s_client -connect localhost:8080 -quiet
```

Example session:

```
PING
PONG
SET name redis
OK
GET name
redis
SET_WITH_TTL session abc123 60
OK (TTL: 60)
GET session
abc123
DELETE name
true
GET name
(nil)
DELETE nonexistent
false
```

### Build

```bash
go build -o redis-server ./cmd/server
./redis-server
```

## TTL and Key Expiration

Keys can be set with a Time-To-Live (TTL) using `SET_WITH_TTL`. Expiration is handled in two layers:

- **Lazy check:** `GET` checks if a key is expired before returning. Expired keys return `(nil)` without blocking other readers.
- **Active cleanup:** A background goroutine runs every second, scanning the store and removing expired keys with a write lock.

This dual approach keeps read performance high (`RLock`) while preventing memory leaks from unaccessed expired keys.

## Concurrency

The store uses `sync.RWMutex` for thread safety:

- **Read operations** (`GET`) use `RLock` - multiple readers can access simultaneously
- **Write operations** (`SET`, `DELETE`, cleanup) use `Lock` - exclusive access, blocks all other reads and writes

Each client connection runs in its own goroutine, so multiple clients can connect and send commands concurrently.

The background cleanup goroutine is started automatically with `store.New()` and can be stopped gracefully with `store.Close()`.

## TLS

The server supports TLS encryption out of the box. When `CertFile` and `KeyFile` are set on `TCPServer`, it uses `tls.Listen` instead of `net.Listen`. If both fields are empty, the server falls back to plain TCP.

```go
// TLS enabled
server := connection.TCPServer{
    Handler:  r.Handle,
    CertFile: "cert.pem",
    KeyFile:  "key.pem",
}

// Plain TCP (no TLS)
server := connection.TCPServer{
    Handler: r.Handle,
}
```
