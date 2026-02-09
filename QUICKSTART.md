# Quick Start Guide

## 1. Setup Project

### Install Go
Pastikan Go versi 1.21+ sudah terinstall:
```bash
go version
```

### Clone/Navigate ke Project
```bash
cd c:\Project-projectan\ProjectWebCurhat
```

### Download Dependencies
```bash
go mod download
```

## 2. Running the Server

### Development Mode
```bash
go run main.go
```

Server akan running di `http://localhost:8080`

### Production Build
```bash
# Build executable
go build -o bin/server.exe .

# Run executable
.\bin\server.exe
```

### Custom Port
```bash
# Windows PowerShell
$env:PORT="3000"
go run main.go

# Windows CMD
set PORT=3000
go run main.go
```

## 3. Testing

### Method 1: Menggunakan Test Client HTML

1. **Start server**:
   ```bash
   go run main.go
   ```

2. **Buka test client**:
   - Buka `test-client.html` di browser Chrome/Firefox
   - Buka tab baru dan buka `test-client.html` lagi (2 tabs)

3. **Test video call**:
   - Tab 1: Enter username "User1", klik "Connect"
   - Tab 2: Enter username "User2", klik "Connect"
   - Video call akan otomatis terkoneksi

### Method 2: Menggunakan 2 Browser Berbeda

1. Start server
2. Buka `test-client.html` di Chrome → Connect
3. Buka `test-client.html` di Firefox → Connect
4. Kedua user akan otomatis matched

### Method 3: Test di 2 Device Berbeda

1. **Start server dengan bind ke semua interface**:
   ```go
   // Tidak perlu ubah apa-apa, sudah default bind ke 0.0.0.0
   ```

2. **Cari IP address komputer**:
   ```bash
   ipconfig
   # Cari IPv4 Address misalnya: 192.168.1.100
   ```

3. **Update WS_URL di test-client.html**:
   ```javascript
   const WS_URL = 'ws://192.168.1.100:8080/ws';
   ```

4. **Buka di 2 device berbeda**:
   - Device 1: Buka `http://192.168.1.100:8080/test-client.html`
   - Device 2: Buka `http://192.168.1.100:8080/test-client.html`

## 4. API Testing dengan WebSocket Client

### Menggunakan wscat (WebSocket CLI)

Install wscat:
```bash
npm install -g wscat
```

Connect ke server:
```bash
wscat -c "ws://localhost:8080/ws?username=TestUser"
```

Send messages:
```json
# Join room
{"type":"join","username":"TestUser1"}

# Leave room
{"type":"leave"}
```

## 5. Health Check

Test server status:
```bash
# PowerShell
Invoke-WebRequest http://localhost:8080/health

# Atau buka di browser
http://localhost:8080/health
```

Response:
```json
{
  "success": true,
  "message": "Server is running",
  "data": {
    "status": "healthy"
  }
}
```

## 6. Monitoring Logs

Saat running, server akan menampilkan logs:

```
2026/02/09 10:00:00 Starting WebRTC signaling server on :8080
2026/02/09 10:00:00 WebSocket endpoint: ws://localhost:8080/ws
2026/02/09 10:00:15 New client connected: abc-123 (username: User1)
2026/02/09 10:00:15 Created new room def-456 for client abc-123
2026/02/09 10:00:20 New client connected: xyz-789 (username: User2)
2026/02/09 10:00:20 Client xyz-789 joined existing room def-456
2026/02/09 10:00:20 Room def-456 is ready with clients abc-123 and xyz-789
```

## 7. Common Issues & Solutions

### Issue: "bind: address already in use"
**Solution**: Port 8080 sudah dipakai, ganti port:
```bash
$env:PORT="3001"
go run main.go
```

### Issue: Camera/Mic tidak terdeteksi di test client
**Solution**: 
- Pastikan browser punya permission untuk camera/mic
- Buka di `https://` atau `localhost` (not `file://`)
- Coba browser lain

### Issue: Video tidak muncul tapi connected
**Solution**:
- Check browser console untuk error
- Pastikan kedua client sudah allow camera/mic
- Check firewall settings

### Issue: Koneksi terputus terus
**Solution**:
- Check STUN server masih accessible
- Pastikan tidak ada proxy/firewall blocking WebSocket
- Try different STUN servers di test-client.html

## 8. Development Tips

### Hot Reload (Auto Restart on File Change)

Install air:
```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main.exe ."
  bin = "tmp/main.exe"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
```

Run with air:
```bash
air
```

### Debugging

Tambahkan debug logs di code:
```go
import "log"

log.Printf("Debug: client %s, room %s", client.ID, room.ID)
```

### Testing dengan curl

```bash
# Health check
curl http://localhost:8080/health

# Home
curl http://localhost:8080/
```

## 9. Production Deployment

### Using systemd (Linux)

Create service file `/etc/systemd/system/webcurhat.service`:
```ini
[Unit]
Description=WebRTC Signaling Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/webcurhat
ExecStart=/opt/webcurhat/server
Restart=always
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable webcurhat
sudo systemctl start webcurhat
```

### Using Docker

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

Build and run:
```bash
docker build -t webcurhat .
docker run -p 8080:8080 webcurhat
```

### Using Windows Service

Gunakan tools seperti:
- NSSM (Non-Sucking Service Manager)
- WinSW

Example dengan NSSM:
```bash
nssm install WebCurhat "C:\path\to\server.exe"
nssm start WebCurhat
```

## 10. Next Steps

- [ ] Add HTTPS/WSS support
- [ ] Implement authentication
- [ ] Add rate limiting
- [ ] Setup reverse proxy (nginx/caddy)
- [ ] Add monitoring (Prometheus/Grafana)
- [ ] Implement turn server untuk NAT traversal
- [ ] Add database untuk chat history
- [ ] Deploy to cloud (AWS/GCP/Azure)

## 11. Resources

- **WebRTC**: https://webrtc.org/
- **Gorilla WebSocket**: https://github.com/gorilla/websocket
- **Go Documentation**: https://go.dev/doc/

## Need Help?

Check:
1. [README.md](README.md) - Overview dan fitur
2. [ARCHITECTURE.md](ARCHITECTURE.md) - Detailed architecture
3. Logs output - Server logs menampilkan detailed info
4. Browser console - Client-side errors

Happy coding! 🚀
