# 📦 DocVault — Document Storage API (Project 1 of 3)

> A production-style REST API for uploading, storing, managing, and auto-expiring documents — the foundation of a complete Document Management System.

## 🎯 System Overview

This is **Project 1** of a 3-project continuous system:

```
┌──────────────────────────────────────────────────────────────────┐
│                  Document Management System                      │
│                                                                  │
│  Project 1: DocVault        ──► Core API (upload, download,      │
│  (You are here)                 store, delete, SQS events)       │
│                                                                  │
│  Project 2: GoAuth          ──► Authentication layer              │
│                                 (JWT, login, RBAC, protects      │
│                                  DocVault routes)                 │
│                                                                  │
│  Project 3: GoFlow          ──► Document processor                │
│                                 (consumes DocVault's SQS events, │
│                                  extracts text, detects dupes,    │
│                                  generates stats)                 │
│                                                                  │
│  Shared: SQLite (users + documents + processing results)         │
│  Shared: MinIO (file storage)                                    │
│  Shared: LocalStack SQS (event bus between DocVault & GoFlow)    │
└──────────────────────────────────────────────────────────────────┘
```

## 🧰 Tech Stack

| Tool | Purpose | Why This? |
|------|---------|-----------|
| **Go (Gin)** | HTTP framework | Lightweight, fast, great docs |
| **MinIO** | Object storage (S3-compatible) | Free, local, identical API to AWS S3 |
| **SQLite** | Database | Zero config, shared across all 3 projects |
| **LocalStack SQS** | Message queue | Event bus between DocVault → GoFlow |
| **Docker Compose** | Container orchestration | Run MinIO + LocalStack together |

---

## 📚 Go Concepts You Will Learn

### Core Language Fundamentals
- [ ] Structs and struct tags (`json:"name"`, `db:"column"`)
- [ ] Methods on structs (value vs pointer receivers)
- [ ] Interfaces (`io.Reader`, `io.Writer`, custom repository/service interfaces)
- [ ] Interface-based design (depend on abstractions, not concretions)
- [ ] Error handling (`value, err` pattern, `fmt.Errorf` wrapping)
- [ ] Package organization (importing your own packages)
- [ ] Defer and cleanup patterns (`defer file.Close()`)
- [ ] Pointers and references
- [ ] Type conversions

### HTTP & Web
- [ ] Gin router setup and route groups
- [ ] HTTP methods (GET, POST, DELETE)
- [ ] Path parameters (`:id`) and query parameters
- [ ] Multipart form-data file uploads
- [ ] Streaming file downloads with proper headers
- [ ] JSON request/response encoding/decoding
- [ ] HTTP status codes (200, 201, 400, 404, 500)

### Database & Persistence
- [ ] `database/sql` interface
- [ ] SQLite driver (`go-sqlite3`)
- [ ] Table creation (migrations)
- [ ] SQL prepared statements
- [ ] CRUD operations
- [ ] Database error handling and connection management

### Concurrency
- [ ] Goroutines (`go func()`)
- [ ] Channels (`<-chan` from MinIO's ListObjects)
- [ ] `select` statement (multiplexing signals)
- [ ] `context.Context` (timeouts, cancellation, propagation)
- [ ] `sync.WaitGroup` (waiting for workers)
- [ ] Background task processing (SQS consumer)

### Clean Architecture & Patterns
- [ ] Entity layer (pure domain structs, zero external deps)
- [ ] Repository pattern (interface + SQLite implementation)
- [ ] Service interfaces (storage, queue abstractions)
- [ ] Usecase layer (business logic orchestration)
- [ ] Handler layer (HTTP parsing, DTO conversion)
- [ ] DTO pattern (separate API contract from domain entity)
- [ ] Factory pattern (centralized dependency injection)
- [ ] Middleware pattern (validation, logging)
- [ ] Dependency inversion

### Unit Testing
- [ ] `testing` package basics
- [ ] Table-driven tests
- [ ] Mock interfaces with testify/mock
- [ ] Testing usecases with mock repository/service
- [ ] Testing handlers with `httptest`
- [ ] `go test -cover`, `go test -race`

### Configuration & Resilience
- [ ] Environment variables and config struct
- [ ] Graceful shutdown (`os.Signal`, `signal.Notify`)
- [ ] Panic recovery middleware
- [ ] Request logging

### External SDKs
- [ ] MinIO Go SDK (`minio-go/v7`)
- [ ] AWS SDK for Go v2 (SQS)
- [ ] Docker Compose

---

## 📁 Project Structure (Clean Architecture + Factory DI)

```
docvault/
├── main.go                     # Entry point: factory → server + workers → graceful shutdown
├── go.mod / go.sum
├── .env
├── docker-compose.yml          # MinIO + LocalStack
│
├── config/
│   └── config.go               # Load env vars into Config struct
│
├── factory/
│   └── factory.go              # Creates ALL dependencies, wires them together
│
├── entity/
│   ├── document.go             # Domain struct — NO JSON tags, NO external deps
│   └── event.go                # Domain event (Type, DocumentID, Filename, Timestamp)
│
├── dto/
│   ├── request.go              # UploadRequest, ListQuery
│   └── response.go             # DocumentResponse, ListResponse, ErrorResponse + FromEntity()
│
├── repository/
│   ├── repository.go           # Interface: DocumentRepository
│   └── sqlite_document.go      # SQLite implementation
│
├── service/
│   ├── storage.go              # Interface: StorageService
│   ├── storage_minio.go        # MinIO implementation
│   ├── queue.go                # Interface: QueueService
│   └── queue_sqs.go            # SQS implementation
│
├── usecase/
│   └── document.go             # Business logic — depends ONLY on interfaces
│
├── handler/
│   └── document.go             # HTTP handlers — depends ONLY on usecase + dto
│
├── database/
│   ├── sqlite.go               # Open connection
│   └── migrations.go           # CREATE TABLE (documents table + processing_results table for GoFlow)
│
├── worker/
│   ├── notification.go         # SQS consumer goroutine
│   └── scheduler.go            # Cron job: auto-delete expired files
│
├── middleware/
│   ├── validation.go           # File type/size validation
│   └── logging.go              # Request logging
│
└── tests/
    ├── mock/
    │   ├── mock_document_repo.go
    │   ├── mock_storage_svc.go
    │   └── mock_queue_svc.go
    ├── usecase/
    │   └── document_test.go
    └── handler/
        └── document_test.go
```

### Clean Architecture Layer Rules

```
┌─────────────────────────────────────────────────────────────┐
│  entity/       ZERO imports. Pure domain. Heart of the app. │
├─────────────────────────────────────────────────────────────┤
│  repository/   Interface + SQLite impl. Depends on: entity  │
├─────────────────────────────────────────────────────────────┤
│  service/      Interface + impl (MinIO, SQS). Depends on:   │
│                entity                                        │
├─────────────────────────────────────────────────────────────┤
│  usecase/      Business logic. Depends on: entity +          │
│                repo/service INTERFACES only                  │
├─────────────────────────────────────────────────────────────┤
│  handler/      HTTP layer. Depends on: usecase + dto only    │
├─────────────────────────────────────────────────────────────┤
│  dto/          JSON tags live here. entity ↔ DTO conversion  │
├─────────────────────────────────────────────────────────────┤
│  factory/      Wires everything. THE ONLY package that       │
│                imports all layers                             │
└─────────────────────────────────────────────────────────────┘
```

### How the layers connect

```
HTTP Request → middleware/ → handler/ → usecase/ → repository/ (SQLite)
                                                 → service/storage (MinIO)
                                                 → service/queue (SQS → GoFlow picks up)

Background:  worker/notification.go → service/queue (consume)
             worker/scheduler.go    → usecase/ (delete expired)

Wiring:      factory/ → creates everything
             main.go  → calls factory, starts server + workers
```

### Factory Pattern

```go
// factory/factory.go
type Factory struct {
    DocumentHandler    *handler.DocumentHandler
    NotificationWorker *worker.NotificationWorker
    SchedulerWorker    *worker.SchedulerWorker
}

func New(cfg *config.Config) (*Factory, error) {
    db, err := database.NewSQLite(cfg.DBPath)            // shared DB
    docRepo := repository.NewSQLiteDocumentRepo(db)
    storageSvc := service.NewMinIOStorage(minioClient, cfg.BucketName)
    queueSvc := service.NewSQSQueue(sqsClient, cfg.QueueURL)
    docUsecase := usecase.NewDocumentUsecase(docRepo, storageSvc, queueSvc)
    docHandler := handler.NewDocumentHandler(docUsecase)
    notifWorker := worker.NewNotificationWorker(queueSvc)
    schedWorker := worker.NewSchedulerWorker(docUsecase)
    return &Factory{docHandler, notifWorker, schedWorker}, nil
}
```

```go
// main.go
func main() {
    cfg := config.Load()
    f, _ := factory.New(cfg)

    r := gin.Default()
    api := r.Group("/api")
    {
        api.POST("/documents/upload", f.DocumentHandler.Upload)
        api.GET("/documents", f.DocumentHandler.List)
        api.GET("/documents/:id", f.DocumentHandler.GetMetadata)
        api.GET("/documents/:id/download", f.DocumentHandler.Download)
        api.DELETE("/documents/:id", f.DocumentHandler.Delete)
        api.GET("/documents/expiring", f.DocumentHandler.ListExpiring)
    }
    r.GET("/health", f.DocumentHandler.Health)

    ctx, cancel := context.WithCancel(context.Background())
    go f.NotificationWorker.Start(ctx)
    go f.SchedulerWorker.Start(ctx)
    // graceful shutdown...
}
```

---

## 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/documents/upload` | Upload file (multipart) → MinIO + SQLite + SQS event |
| `GET` | `/api/documents` | List all docs (filter by type, sort) |
| `GET` | `/api/documents/:id` | Get file metadata |
| `GET` | `/api/documents/:id/download` | Stream file download |
| `DELETE` | `/api/documents/:id` | Delete from MinIO + SQLite + SQS event |
| `GET` | `/api/documents/expiring?within=7` | List files expiring soon |
| `GET` | `/health` | Health check (SQLite + MinIO + SQS) |

> **Note:** These routes are currently unprotected. In Project 2 (GoAuth), you'll add JWT authentication middleware to protect them.

---

## 🗺️ Phase-by-Phase Roadmap

### Phase 1: Project Setup & Factory Foundation (Day 1–2)

**Goal:** Docker services, Go project, database, entity, factory skeleton, health check.

**Steps:**
1. `go mod init docvault`
2. Write `docker-compose.yml` (MinIO + LocalStack)
3. Create `config/config.go`, `database/sqlite.go`, `database/migrations.go`
4. Create `entity/document.go` (ZERO imports), `entity/event.go`
5. Create `factory/factory.go` (skeleton), `main.go`
6. Add `GET /health` endpoint

**Important for continuity:** In `database/migrations.go`, also create the `users` table and `processing_results` table now — GoAuth and GoFlow will use them later:
```sql
CREATE TABLE IF NOT EXISTS documents (...);
CREATE TABLE IF NOT EXISTS users (...);           -- GoAuth will use this
CREATE TABLE IF NOT EXISTS processing_results (...); -- GoFlow will write here
```

**Test yourself:**
- [ ] `docker-compose up -d` starts MinIO and LocalStack
- [ ] `go run main.go` starts the server
- [ ] `curl http://localhost:8080/health` returns OK
- [ ] `entity/document.go` has ZERO external imports

---

### Phase 2: Upload — Entity → Repo → Service → Usecase → Handler (Day 3–4)

**Goal:** File upload following clean architecture, bottom-up.

**Steps:**
1. `repository/repository.go` — `DocumentRepository` interface + `sqlite_document.go` impl
2. `service/storage.go` — `StorageService` interface + `storage_minio.go` impl
3. `usecase/document.go` — `DocumentUsecase` with Upload logic
4. `dto/response.go` — `DocumentResponse` + `FromEntity()` conversion
5. `handler/document.go` — Upload handler
6. `middleware/validation.go` — file type/size check
7. Wire in `factory/factory.go`, register route in `main.go`

**Test yourself:**
- [ ] Upload works: `curl -F "file=@test.pdf" http://localhost:8080/api/documents/upload`
- [ ] File in MinIO, metadata in SQLite
- [ ] `usecase/` does NOT import `minio-go` or `database/sql`
- [ ] `handler/` does NOT import `repository/` or `service/`

---

### Phase 3: List & Download (Day 5)

**Goal:** List documents from SQLite, stream downloads from MinIO.

**Steps:** Add FindAll, FindByID to repo → Download to storage service → List, GetMetadata, Download to usecase → handlers → wire in factory

**Test yourself:**
- [ ] `curl http://localhost:8080/api/documents` returns list
- [ ] `curl -o output.pdf http://localhost:8080/api/documents/:id/download` streams file
- [ ] Uses `io.Copy` to stream (not loading into memory)

---

### Phase 4: Delete (Day 6)

**Goal:** Delete from both MinIO and SQLite via usecase orchestration.

**Steps:** Add Delete to repo + storage service → DeleteDocument usecase → Delete handler

**Test yourself:**
- [ ] File gone from MinIO + SQLite after delete
- [ ] Deletion logic in usecase, NOT handler

---

### Phase 5: SQS Integration (Day 7)

**Goal:** Publish events to SQS on upload/delete — GoFlow will consume these in Project 3.

**Steps:**
1. `service/queue.go` — `QueueService` interface (Publish, Consume, DeleteMessage)
2. `service/queue_sqs.go` — LocalStack SQS implementation
3. Update usecase — publish `file.uploaded` / `file.deleted` events AFTER success
4. Update factory

**Events published (GoFlow consumes these later):**
```json
{"type": "file.uploaded", "document_id": "...", "filename": "report.pdf", "content_type": "application/pdf", "timestamp": "..."}
{"type": "file.deleted", "document_id": "...", "filename": "report.pdf", "timestamp": "..."}
```

**Test yourself:**
- [ ] Upload → SQS has `file.uploaded` message
- [ ] Delete → SQS has `file.deleted` message
- [ ] `usecase/` uses `QueueService` interface, NOT `aws-sdk-go`

---

### Phase 6: Notification Worker (Day 8)

**Goal:** Background goroutine consuming SQS messages.

**Steps:** Create `worker/notification.go` → Start as goroutine in `main.go` → Context cancellation

**Test yourself:**
- [ ] Upload → notification log appears in terminal
- [ ] Ctrl+C → worker stops gracefully

---

### Phase 7: Scheduler — Auto-Delete Expired (Day 9–10)

**Goal:** Background cron job to find and delete expired documents.

**Steps:** Add FindExpired to repo → DeleteExpiredDocuments to usecase → `worker/scheduler.go` with `time.Ticker` → ListExpiring handler for API

**Test yourself:**
- [ ] File with `expires_in=0` auto-deleted by scheduler
- [ ] SQS has `file.expired` event
- [ ] Scheduler calls usecase, NOT repo/service directly

---

### Phase 8: Unit Tests (Day 11–12)

**Goal:** Test usecase and handler layers with mocks.

**Steps:**
1. Create mocks for DocumentRepository, StorageService, QueueService
2. Test usecase: Upload success/failure, Delete success/not-found, DeleteExpired batch
3. Test handlers with `httptest.NewRecorder`
4. Table-driven tests for multiple scenarios

**Test yourself:**
- [ ] `go test ./...` passes, `go test -race ./...` clean
- [ ] `go test -cover ./...` shows coverage
- [ ] ZERO real MinIO/SQLite calls in tests

---

### Phase 9: Polish & Hardening (Day 13–14)

**Goal:** Graceful shutdown, logging, panic recovery, architecture verification.

**Steps:** Signal handling + WaitGroup → logging middleware → panic recovery → health check pings all services → verify import rules

**Test yourself:**
- [ ] Ctrl+C → clean shutdown
- [ ] Consistent error format: `{"error": "message"}`
- [ ] No import cycle violations

---

## 🧪 Testing Cheat Sheet

```bash
docker-compose up -d
go run main.go

curl -F "file=@test.pdf" -F "expires_in=30" http://localhost:8080/api/documents/upload
curl http://localhost:8080/api/documents
curl http://localhost:8080/api/documents/<id>
curl -o output.pdf http://localhost:8080/api/documents/<id>/download
curl -X DELETE http://localhost:8080/api/documents/<id>
curl http://localhost:8080/api/documents/expiring?within=7
curl http://localhost:8080/health

# SQS check
aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/docvault-events

# Tests
go test ./...
go test -cover ./...
go test -race ./...
```

---

## 📖 Key Dependencies

```bash
go get github.com/gin-gonic/gin
go get github.com/minio/minio-go/v7
go get github.com/mattn/go-sqlite3
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/service/sqs
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/google/uuid
go get github.com/stretchr/testify
```

---

## 🐳 Docker Compose

```yaml
version: '3.8'
services:
  minio:
    image: minio/minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

  localstack:
    image: localstack/localstack
    ports:
      - "4566:4566"
    environment:
      SERVICES: sqs
      DEFAULT_REGION: us-east-1

volumes:
  minio_data:
```

---

## 💡 Tutor Instructions

1. **Don't give full solutions** — give function signatures, let me implement
2. **Point to package docs** — "look at minio-go's PutObject method"
3. **Ask me to explain errors** before fixing
4. **Check my architecture** — handler must not import repository
5. **Ask me why** — "why does usecase take an interface?"

### Common mistakes to watch for:
- Not closing file handles (use `defer`)
- Loading files into memory instead of streaming
- JSON tags on entity structs (belong in DTOs)
- Handler importing repository (should go through usecase)
- Usecase importing `minio-go` (should use interface)
- Business logic in handler (move to usecase)
- Forgetting to update factory when adding dependencies
- SQL injection (use prepared statements)
- Not deleting SQS messages after processing

---

## ✅ Completion Checklist

- [ ] Phase 1: Setup, Docker, SQLite (all 3 tables), entity, factory skeleton
- [ ] Phase 2: Upload (entity → repo → service → usecase → handler → factory)
- [ ] Phase 3: List + download
- [ ] Phase 4: Delete via usecase
- [ ] Phase 5: SQS events (GoFlow will consume these!)
- [ ] Phase 6: Notification worker
- [ ] Phase 7: Scheduler (auto-delete expired)
- [ ] Phase 8: Unit tests with mocks
- [ ] Phase 9: Graceful shutdown, logging, polish

**→ Next: Project 2 (GoAuth) — Add authentication to protect these routes**