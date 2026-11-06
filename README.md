# Layrd Framework (Go Universal Backend Chassis)

> **Layrd** (/ˈleɪ.əd/) adalah starter framework backend modular berbasis **Golang** untuk membangun REST API dengan arsitektur berlapis dan module yang dapat diperluas. Evaluasi performa dan kesiapan production bergantung pada aplikasi yang dibangun di atasnya.

---

## 📦 Go Framework
```bash
git clone github.com/rafialfauzan/layrd
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
6. **Built-in Middlewares**
   Termasuk CORS dinamis, JWT Authentication & RBAC role checker, IP/Token Token-Bucket Rate Limiter, Panic Recovery dengan stack tracing aman, dan Unique `X-Request-ID`.
7. **Graceful Shutdown**  
   Mencegah hilangnya koneksi saat server restart atau deploy dengan menangani sinyal OS `SIGINT` / `SIGTERM` secara tertib.

---

## 📁 Struktur Direktori

```
barong-framework/
├── cmd/
│   ├── server/
│   │   └── main.go                  # Entry point aplikasi
│   └── migrate/
│       └── main.go                  # Perintah migrasi schema eksplisit
├── internal/
│   ├── bootstrap/                   # Inisialisasi aplikasi dan wiring
│   │   ├── bootstrap.go             # Orkestrasi bootstrap dan server
│   │   ├── dependencies.go          # Registrasi dependency/module
│   │   ├── routes.go                # Registrasi route aplikasi
│   │   └── storage.go               # Setup static file storage
│   ├── buildinfo/                   # Metadata versi yang diisi saat build
│   └── modules/                     # Domain modules aplikasi
│       └── user/                    # Modul User (model, repo, service, handler, routes)
│           ├── model.go
│           ├── repository.go        # Embed BaseRepository[User]
│           ├── service.go
│           ├── handler.go
│           └── module.go
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

### 3. Daftarkan Modul di `internal/bootstrap/dependencies.go`
```go
return app.RegisterModule(product.NewModule(app))
```

### Keamanan dan migrasi

Semua endpoint pada module user membutuhkan `Authorization: Bearer <JWT>`.
Layrd belum menyediakan endpoint login/token; token harus diterbitkan oleh mekanisme
autentikasi aplikasi. Role dan status aktif user tidak dapat diubah lewat request
module user, dan user baru selalu mendapat role `user`.
Semua token valid saat ini memiliki akses yang sama ke operasi CRUD user; tambahkan
otorisasi berbasis kepemilikan atau role sebelum aplikasi multi-user digunakan.

Schema tidak dimigrasikan otomatis saat server start. Setelah mengatur koneksi
database dan membuat database kosong, jalankan migrasi secara eksplisit:
```bash
go run ./cmd/migrate
```
Perintah ini menjalankan migrator module yang terdaftar. Tinjau perubahan schema
dan gunakan mekanisme backup/roll-forward yang sesuai sebelum menjalankannya di production.

---

## 🗄️ Mengakses Multi-Database

Layrd mendukung PostgreSQL (`DB_DRIVER=postgres`, default) dan MySQL (`DB_DRIVER=mysql`).
Driver dapat dipilih terpisah untuk database utama dan database sekunder melalui
`DB_USER_DRIVER`. Port default mengikuti driver: `5432` untuk PostgreSQL dan `3306`
untuk MySQL.

Contoh konfigurasi database utama MySQL di `.env`:
```dotenv
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=layrd_main
```

Database sekunder mewarisi driver utama secara default; atur `DB_USER_DRIVER`
untuk memilih driver yang berbeda.

Di environment production, `JWT_SECRET` harus ditentukan eksplisit dan memiliki
panjang minimal 32 karakter.

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
