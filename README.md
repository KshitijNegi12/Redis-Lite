# Redis Lite Project

A custom lightweight Redis-like server implementation written in Go, featuring a React-based Admin Panel for interacting with the server.

This project implements core Redis functionality, including standard data manipulations, key expiries, transactions, streams, and master-replica architectures. A companion Web Admin Panel allows you to easily manage and send commands directly to the server.

## Features

- **Core Redis Commands**: Supports essential commands such as `SET`, `GET`, `DEL`, `ECHO`, `PING`, and `INCR`.
- **Key Expiry**: Fully supports `EX` (Seconds) and `PX` (Milliseconds) options for the `SET` command.
- **Master-Replica Replication**: Supports replication concepts using `REPLCONF`, `PSYNC`, and configurable start-up options (`--replicaof`).
- **Transactions**: Supports `MULTI`, `EXEC`, and `DISCARD` for batch command execution.
- **Streams**: Basic support for `XADD` and `XREAD` streams.
- **Admin Dashboard**: A React-based web panel that interfaces with the server through a custom Node.js WebSocket proxy, providing an interactive console to send commands.

---

## Project Structure

- **`RedisServer/`**: The core Redis engine written in Go.
  - Implements the RESP (REdis Serialization Protocol).
  - Handles command parsing (`RDBparser`, `resp/`).
  - Contains command executions (`implementation/`, `handler/`).
  - Manages connections and server loop (`server/`).
- **`Admin/`**: The React Admin Frontend.
  - A real-time monitoring and interactive web console (`src/App.js`).
  - **`Admin/src/proxy/proxyServer.js`**: A Node.js middleware server that bridges WebSocket connections from the Admin web panel to the TCP-based native Go Redis server.

---

## Supported Commands

| Command        | Description                                                                                                                 |
| -------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **`PING`**     | Returns `PONG`. Useful for testing connections.                                                                             |
| **`ECHO`**     | Returns the given string/message.                                                                                           |
| **`SET`**      | Set `key` to hold the `value`. Options: `EX` (seconds) and `PX` (milliseconds) for expiry. Example: `SET mykey value EX 60` |
| **`GET`**      | Get the value of `key`. Returns nil if the key does not exist or has expired.                                               |
| **`DEL`**      | Deletes the specified `key`.                                                                                                |
| **`INFO`**     | Returns server information and statistics (useful for replication info).                                                    |
| **`KEYS`**     | Returns all keys matching `pattern`. Example: `KEYS *` or `KEYS mykey?`                                                     |
| **`REPLCONF`** | Used by replicas to configure the replication process.                                                                      |
| **`PSYNC`**    | Used by replicas to synchronize with the master.                                                                            |
| **`TYPE`**     | Returns the type of the value stored at `key`.                                                                              |
| **`XADD`**     | Appends a new entry to a stream.                                                                                            |
| **`XREAD`**    | Reads data from one or multiple streams.                                                                                    |
| **`INCR`**     | Increments the number stored at `key` by one.                                                                               |
| **`MULTI`**    | Marks the start of a transaction block.                                                                                     |
| **`EXEC`**     | Executes all commands within a transaction block (`MULTI`).                                                                 |
| **`DISCARD`**  | Flushes/cancels all previously queued commands in a transaction block.                                                      |

---

## How to Run the Project

To run the full stack, you need to start the backend Redis server, the proxy engine, and the frontend admin panel on your local machine. Ensure you have Go (1.23+) and Node.js installed.

### 1. Start the Go Redis Server

Navigate to the `RedisServer` directory and run the Go application. By default, it runs on port `6379` as a master node.

```bash
cd RedisServer
go run main.go
```

**Optional Command-Line Arguments**:

- `--port <port>`: Starts the server on a custom port.
- `--replicaof <master-host>:<master-port>`: Starts the server as a replica node of the specified master.
- `--master_replid <id>`: Specifies a custom master replication ID.

_Example (Starting a replica on port 6380):_

```bash
go run main.go --port 6380 --replicaof 127.0.0.1:6379
```

### 2. Start the WebSocket Proxy Server

The Admin panel requires a bridge to communicate with raw TCP over WebSockets. Run the Node.js proxy server.

```bash
cd Admin/src/proxy
node proxyServer.js
```

_This proxy will run on `localhost:8080` by default and route WebSocket messages to RESP commands on Redis port 6379 (configurable in the Admin web UI)._

### 3. Start the React Admin Panel

In a new terminal, run the React frontend application.

```bash
cd Admin
npm install
npm start
```

_The web dashboard should automatically open in your default browser at `http://localhost:3000`. In the web UI, ensure the connection port matches your chosen Redis server port (e.g., 6379), and start typing commands in the real-time interactive terminal!_
