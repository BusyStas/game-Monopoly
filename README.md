# game-Monopoly

A prototype, single-binary multiplayer board game server + embedded browser UI.

## What's included

- Minimal Go HTTP server (serves an embedded static UI) with a simple WebSocket hub.

## Build & run

Build:

```bash
go build -o bin/game-monopoly ./...
```

Run (dev):

```bash
go run .
# then open http://localhost:8080
```

The server serves the embedded static UI and provides a websocket endpoint at `/ws`.

## Next steps

- Add UDP peer discovery and LAN messaging.
- Implement full game logic, player state, and front-end board UI.
