# ProjectWebCurhat - WebRTC P2P Video Call

Backend signaling server untuk aplikasi video call peer-to-peer menggunakan WebRTC dan Golang.

## Arsitektur

Aplikasi ini menggunakan arsitektur **Clean Architecture** dengan layer-layer berikut:

```
ProjectWebCurhat/
├── main.go                          # Entry point aplikasi
├── go.mod                           # Dependencies
├── config/
│   └── config.go                    # Configuration management
├── internal/
│   ├── model/                       # Layer Model (Domain Entities)
│   │   ├── client.go               # WebSocket client model
│   │   ├── message.go              # WebRTC signaling messages
│   │   └── room.go                 # Chat room model
│   ├── service/                     # Layer Business Logic
│   │   ├── room_service.go         # Room management service
│   │   └── signaling_service.go    # WebRTC signaling service
│   └── handler/                     # Layer Presentation (HTTP/WebSocket)
│       └── websocket_handler.go    # WebSocket connection handler
└── pkg/
    └── response/                    # Layer Utilities
        └── response.go             # HTTP response helper
```

## Penjelasan Layer

### 1. **Model Layer** (`internal/model/`)

- Berisi domain entities dan struktur data
- `Client`: Representasi WebSocket client
- `Room`: Representasi ruang chat untuk 2 peer
- `Message`: Struktur pesan signaling WebRTC

### 2. **Service Layer** (`internal/service/`)

- Berisi business logic aplikasi
- `RoomService`: Mengelola room matching dan lifecycle
- `SignalingService`: Menangani WebRTC signaling (offer/answer/ICE)

### 3. **Handler Layer** (`internal/handler/`)

- Menangani HTTP/WebSocket requests
- `WebSocketHandler`: Handle koneksi WebSocket dan routing pesan

### 4. **Config Layer** (`config/`)

- Mengelola konfigurasi aplikasi
- Load dari environment variables

### 5. **Utilities** (`pkg/`)

- Helper functions yang reusable
- Response formatting

## Cara Menjalankan

### 1. Install Dependencies

```bash
go mod download
```

### 2. Jalankan Server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

### 3. Environment Variables (Optional)

```bash
# Set port (default: 8080)
$env:PORT="3000"
go run main.go
```

## Endpoints

- **WebSocket**: `ws://localhost:8080/ws?username=YourName`
- **Health Check**: `http://localhost:8080/health`
- **Root**: `http://localhost:8080/`

## WebRTC Signaling Flow

1. **Koneksi**: Client connect ke `/ws` endpoint
2. **Join**: Client kirim message `{"type":"join", "username":"User1"}`
3. **Ready**: Server kirim message `{"type":"ready", "roomId":"..."}`
4. **Matching**: Ketika 2 clients dalam room, server notify keduanya
5. **Offer**: Client pertama kirim SDP offer
6. **Answer**: Client kedua kirim SDP answer
7. **ICE Candidates**: Exchange ICE candidates untuk koneksi
8. **P2P Connection**: Setelah selesai, video/audio stream langsung peer-to-peer

## Format Pesan

### Join Room

```json
{
    "type": "join",
    "username": "User123"
}
```

### SDP Offer/Answer

```json
{
    "type": "offer",
    "payload": {
        "type": "offer",
        "sdp": "v=0\r\no=..."
    }
}
```

### ICE Candidate

```json
{
    "type": "candidate",
    "payload": {
        "candidate": "candidate:...",
        "sdpMid": "0",
        "sdpMLineIndex": 0
    }
}
```

### Leave Room

```json
{
    "type": "leave"
}
```

## Fitur

- ✅ **Peer-to-Peer Matching**: Automatic matching 2 users
- ✅ **WebRTC Signaling**: Handle SDP offer/answer dan ICE candidates
- ✅ **Room Management**: Auto create dan cleanup rooms
- ✅ **Clean Architecture**: Separation of concerns dengan layer yang jelas
- ✅ **Concurrent Safe**: Menggunakan mutex untuk thread safety
- ✅ **WebSocket**: Real-time bidirectional communication

## Testing

Buka file `test-client.html` di 2 browser berbeda untuk test video call.

## Production Checklist

Sebelum deploy ke production:

- [ ] Update CORS configuration di `websocket_handler.go`
- [ ] Tambahkan rate limiting
- [ ] Implementasi authentication/authorization
- [ ] Setup HTTPS/WSS (WebSocket Secure)
- [ ] Add logging dengan level (debug, info, error)
- [ ] Implement graceful shutdown
- [ ] Add monitoring dan metrics
- [ ] Database untuk persistent storage (jika diperlukan)
- [ ] Load balancing configuration
- [ ] Deploy dengan Docker/Kubernetes

## Teknologi

- **Backend**: Go 1.21+
- **WebSocket**: gorilla/websocket
- **UUID**: google/uuid
- **WebRTC**: Browser native WebRTC API

## License

MIT
