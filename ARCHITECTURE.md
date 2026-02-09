# Panduan Arsitektur Berlayer

## Clean Architecture Overview

Aplikasi ini mengimplementasikan **Layered Architecture** dengan pattern **Contract-based Dependency Injection**,
mengikuti pola yang sama dengan `evt-backend-web`. Setiap layer memiliki tanggung jawab yang jelas dan
dependencies hanya mengalir satu arah (inward).

## Struktur Folder

```
main.go                         # Entry point
config/
  config.go                     # App configuration (singleton)
  middleware/
    cors.go                     # CORS middleware (Gin)
  pkg/
    errs/
      errs.go                   # Structured error types
  server/
    server.go                   # Server startup & route initialization
database/
  models.go                     # Domain models (Client, Room, MessageType)
dto/
  websocket.go                  # Data Transfer Objects (message DTOs)
contract/
  repository.go                 # Repository interfaces
  service.go                    # Service interfaces
repository/
  repository.go                 # Repository factory (New)
  room.go                       # Room repository implementation
service/
  service.go                    # Service factory (New)
  room.go                       # Room service implementation
  signaling.go                  # Signaling service implementation
controller/
  controller.go                 # Controller interface & factory (New)
  health.go                     # Health check controller
  websocket.go                  # WebSocket controller
pkg/
  response/
    response.go                 # HTTP response helpers
```

## Layer-Layer dalam Aplikasi

### 1. **Config Layer**
**Lokasi**: `config/`

**Tanggung jawab**:
- Load konfigurasi dari environment variables
- Setup server (Gin engine, middleware, routes)
- CORS middleware
- Structured error types

**File**:
- `config.go`: Singleton config dengan `Load()` dan `Get()`
- `server/server.go`: Inisialisasi server, dependency injection, start HTTP server
- `middleware/cors.go`: CORS middleware untuk Gin
- `pkg/errs/errs.go`: Error types (BadRequest, NotFound, InternalServerError, dll)

---

### 2. **Database Layer** (Domain Models)
**Lokasi**: `database/`

**Tanggung jawab**:
- Mendefinisikan struktur data core (entities)
- Pure Go structs dan methods
- Tidak ada dependencies ke layer lain (selain library)

**File**:
- `models.go`: `Client`, `Room`, `MessageType` - entities utama aplikasi

---

### 3. **DTO Layer** (Data Transfer Objects)
**Lokasi**: `dto/`

**Tanggung jawab**:
- Mendefinisikan format data yang dikirim/diterima via WebSocket
- Memisahkan transport data dari domain models
- JSON serialization tags

**File**:
- `websocket.go`: `Message`, `SDPMessage`, `ICECandidateMessage`, message type constants

---

### 4. **Contract Layer** (Interfaces)
**Lokasi**: `contract/`

**Tanggung jawab**:
- Mendefinisikan interface untuk Repository dan Service
- Menjadi "kontrak" yang harus dipenuhi oleh implementasi
- Memungkinkan loose coupling antar layer

**File**:
- `repository.go`: `Repository` struct + `RoomRepository` interface
- `service.go`: `Service` struct + `RoomService`, `SignalingService` interfaces

---

### 5. **Repository Layer** (Data Access)
**Lokasi**: `repository/`

**Tanggung jawab**:
- Implementasi data access (in-memory room storage)
- CRUD operations untuk Room
- Mengelola waiting room state
- Thread-safe dengan mutex

**File**:
- `repository.go`: Factory function `New()` → `*contract.Repository`
- `room.go`: Implementasi `contract.RoomRepository`

---

### 6. **Service Layer** (Business Logic)
**Lokasi**: `service/`

**Tanggung jawab**:
- Implementasi business logic
- Room matching (find or create)
- WebRTC signaling flow (join, leave, relay)
- Menggunakan Repository untuk data access

**File**:
- `service.go`: Factory function `New(repo) → *contract.Service`
- `room.go`: Implementasi `contract.RoomService`
- `signaling.go`: Implementasi `contract.SignalingService`

---

### 7. **Controller Layer** (Presentation)
**Lokasi**: `controller/`

**Tanggung jawab**:
- Handle HTTP/WebSocket requests
- Upgrade connections ke WebSocket
- Memanggil Service layer untuk business logic
- Setiap controller implements `Controller` interface

**File**:
- `controller.go`: `Controller` interface (`GetPrefix`, `InitService`, `InitRoute`) + factory `New()`
- `health.go`: `HealthController` - endpoint `/health`
- `websocket.go`: `WebSocketController` - endpoint `/ws`

---

## Data Flow

```
Client (WebSocket)
    → Controller Layer (controller/websocket.go)
    → Service Layer (service/signaling.go)
    → Service Layer (service/room.go)
    → Repository Layer (repository/room.go)
    → Database Layer (database/models.go)
    → Repository Layer (return data)
    → Service Layer (process & prepare response)
    → Controller Layer (send response via DTO)
    → Client (WebSocket)
```

## Dependency Flow

```
main.go
  → config.Load()
  → server.Run()
      → repository.New()        → contract.Repository
      → service.New(repo)       → contract.Service
      → controller.New(r, svc)  → registers all routes

Controller ──→ Service (via contract.Service interface)
Service    ──→ Repository (via contract.Repository interface)
Repository ──→ Database (domain models)
```

## Controller Interface Pattern

Setiap controller mengimplementasikan interface:

```go
type Controller interface {
    GetPrefix() string                        // Route prefix (e.g., "/ws")
    InitService(service *contract.Service)    // Inject service dependency
    InitRoute(app *gin.RouterGroup)           // Register routes
}
```

Factory function `New()` iterates semua controller, inject service, dan register routes:

```go
func New(app *gin.Engine, service *contract.Service) {
    allController := []Controller{
        &HealthController{},
        &WebSocketController{},
    }
    for _, c := range allController {
        c.InitService(service)
        group := app.Group(c.GetPrefix())
        c.InitRoute(group)
    }
}
```

## Dependency Rule

```
Controller ──→ Contract(Service) ──→ Contract(Repository) ──→ Database(Models)
     ↑              ↑                       ↑
  Gin HTTP    Service impl            Repository impl
```

Dependencies only point inward. Database models must NOT import Service, Repository, or Controller.


### Menambah Database Layer

Untuk menambah persistence (database):

#### 1. Buat Repository Interface di Service Layer
```go
// internal/service/room_repository.go
type RoomRepository interface {
    Save(room *model.Room) error
    FindByID(id string) (*model.Room, error)
    Delete(id string) error
}
```

#### 2. Implement di Infrastructure Layer
```go
// internal/repository/postgres_room_repository.go
type PostgresRoomRepository struct {
    db *sql.DB
}

func (r *PostgresRoomRepository) Save(room *model.Room) error {
    // Database logic
}
```

#### 3. Inject ke Service
```go
// main.go
db := connectDatabase()
roomRepo := repository.NewPostgresRoomRepository(db)
roomService := service.NewRoomService(roomRepo)
```

## Best Practices

1. **Keep models pure**: Model hanya struct dan basic methods
2. **Fat services**: Business logic ada di service layer
3. **Thin handlers**: Handler hanya routing dan response formatting
4. **Dependency Injection**: Pass dependencies via constructor
5. **Interface segregation**: Gunakan interface untuk dependencies
6. **Single Responsibility**: Setiap file/function satu tanggung jawab

## Testing Strategy

### Unit Tests
- Test setiap layer independently
- Mock dependencies

### Integration Tests
- Test flow antar layer
- Use test database jika ada

### E2E Tests
- Test dari client ke server
- Gunakan test client

## Kesimpulan

Arsitektur berlayer ini membuat aplikasi:
- ✅ Mudah di-maintain
- ✅ Mudah di-test
- ✅ Mudah di-scale
- ✅ Mudah di-extend
- ✅ Clean dan organized
