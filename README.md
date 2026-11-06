# Layrd Framework (Go Universal Backend Chassis)

> **Layrd** (/ˈleɪ.əd/) adalah framework backend modular dan berperforma tinggi berbasis **Golang**, dirancang khusus untuk membangun REST API enterprise dengan arsitektur berlapis (*layered architecture*), bersih, terisolasi, dan mudah diperluas (*pluggable*).

---

## 📦 Go Module
```bash
go get github.com/rafialfauzan/layrd
```

---

## 🌟 Fitur Utama Framework

1. **Modular Architecture (`core.Module`)**  
   Setiap domain fitur (User, Produk, Billing, KDS, Order, dll.) berdiri sendiri sebagai modul independen dengan rute, migrasi, dan dependency-nya sendiri tanpa mengotori `main.go`.
2. **Generic CRUD Repository (`database.BaseRepository[T]`)**  
   Memanfaatkan **Go Generics** sehingga Anda tidak perlu menulis ulang query CRUD (`FindByID`, `FindAll` + pagination/filter, `Create`, `Update`, `Delete`) untuk puluhan entitas model.
3. **Multi-Database Connection Manager (`database.DBManager`)**  
   Dukungan koneksi ke banyak database sekaligus (contoh: Database Resto, Database User, Database Reporting/Trip) dengan manajemen *connection pooling* otomatis.
4. **Real-time WebSocket Hub (`websocket.Hub`)**  
   Sistem WebSocket berbasis *Room/Channel* bawaan yang siap pakai untuk aplikasi kasir/KDS (Kitchen Display System), status pesanan meja, atau notifikasi push.
5. **Standardized Response Envelope (`response.Response`)**  
   Format respon JSON konsisten untuk seluruh endpoint (`Success`, `Paginated`, `Error`, `BadRequest`, `Unauthorized`, `NotFound`, `InternalError`).
6. **Built-in Enterprise Middlewares**  
   Termasuk CORS dinamis, JWT Authentication & RBAC role checker, IP/Token Token-Bucket Rate Limiter, Panic Recovery dengan stack tracing aman, dan Unique `X-Request-ID`.
7. **Graceful Shutdown**  
   Mencegah hilangnya koneksi saat server restart atau deploy dengan menangani sinyal OS `SIGINT` / `SIGTERM` secara tertib.

---

## 📁 Struktur Direktori

```
barong-framework/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point aplikasi
├── pkg/
│   ├── config/
│   │   └── config.go                # Loader konfigurasi lingkungan (.env / OS env)
│   ├── core/
│   │   ├── app.go                   # Orchestrator inti (Router, DB, WebSocket, Server)
│   │   └── module.go                # Kontrak interface modul & lifecycle hooks
│   ├── database/
│   │   ├── manager.go               # Multi-DB manager & connection pooling
│   │   └── generic_repository.go    # Generic Base Repository CRUD (Go Generics)
│   ├── middleware/
│   │   ├── cors.go                  # Dynamic CORS origin handler
│   │   ├── jwt_auth.go              # JWT Validator & Role Access Control
│   │   ├── rate_limiter.go          # IP-based rate limiting (token-bucket)
│   │   ├── recovery.go              # Panic recovery mengembalikan JSON terstruktur
│   │   └── request_id.go            # Request correlation ID generator
│   ├── response/
│   │   └── response.go              # Standard response helper & pagination calculator
│   ├── utils/
│   │   ├── password.go              # Bcrypt hash & password verification
│   │   ├── token.go                 # JWT generator & claim parser
│   │   └── uploader.go              # Helper upload file & validasi gambar/dokumen
│   └── websocket/
│       ├── hub.go                   # Thread-safe multi-room WebSocket hub
│       └── client.go                # Client connection pump & ping/pong
├── internal/
│   ├── assets/                      # Aset statis (logo, stempel, font)
│   ├── modules/                     # Domain modules aplikasi
│   │   └── user/                    # Modul User (model, repo, service, handler, routes)
│   │       ├── model.go
│   │       ├── repository.go        # Embed BaseRepository[User]
│   │       ├── service.go
│   │       ├── handler.go
│   │       └── module.go
│   └── templates/                   # Template HTML untuk email / invoice
├── storage/
│   └── uploads/                     # Tempat file fisik upload pengguna (avatars, dll)
├── .env.example                     # Template konfigurasi environment variables
├── .gitignore                       # Daftar file/folder yang diabaikan oleh Git (bin/, .env, uploads/)
├── build.sh                         # Script kompilasi multi-platform (Linux & Windows) dengan version injection
├── go.mod                           # Definisi modul Go dan dependensi
├── go.sum                           # Checksum verifikasi integritas dependensi
└── README.md                        # Dokumentasi panduan framework dan panduan penggunaan
```

---

## 🚀 Cara Menambah Modul Baru

Membuat modul baru sangat sederhana:

### 1. Definisikan Model
```go
package product

type Product struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `json:"name"`
    Price int64  `json:"price"`
}
```

### 2. Buat Repository dengan Generic Base
```go
package product

import "github.com/rafialfauzan/layrd/pkg/database"

type Repository interface {
    database.BaseRepository[Product]
}

func NewRepository(db *gorm.DB) Repository {
    return &struct{ database.BaseRepository[Product] }{
        BaseRepository: database.NewBaseRepository[Product](db),
    }
}
```
*Anda otomatis mendapatkan `FindByID`, `FindAll` (pagination), `Create`, `Update`, dan `Delete` tanpa coding tambahan.*

### 3. Daftarkan Modul di `cmd/server/main.go`
```go
app.RegisterModule(product.NewModule(app))
```

---

## 🗄️ Mengakses Multi-Database

Jika aplikasi Anda menggunakan lebih dari 1 database (misal `default` dan `user_db`):
```go
// Mengambil database default (Resto)
dbResto := app.DBManager.Default()

// Mengambil database sekunder
dbUser, err := app.DBManager.Get("user")
```

---

## 📡 Menggunakan WebSocket Real-time

Untuk broadcast ke room tertentu (misalnya pesanan meja baru ke Kitchen Display System):
```go
// Broadcast pesan JSON ke room "kitchen-outlet-1"
app.WSHub.BroadcastToRoom("kitchen-outlet-1", orderPayload)
```
Client frontend dapat terhubung ke:
```
ws://localhost:8080/ws?room=kitchen-outlet-1
```
