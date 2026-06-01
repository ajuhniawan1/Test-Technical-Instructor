# Bootcamp Assignment & Submission Platform

Backend API menggunakan **Golang Gin + MySQL + Redis** untuk sistem assignment bootcamp. Project ini dibuat untuk kebutuhan assessment/interview backend dengan fokus pada desain API, authentication, authorization, caching, concurrency handling, dan security improvement.

## Ringkasan Project

Aplikasi ini mengatur proses pembelajaran berbasis class/batch:

1. **Admin** membuat class, assign trainer, dan assign talent.
2. **Trainer** membuat assignment, melihat submission, memberi review, dan meminta revisi.
3. **Talent** melihat assignment, submit tugas, resubmit jika diminta revisi, dan melihat progress.

Backend menggunakan layered architecture:

```text
Request API
  ↓
Gin Router
  ↓
Middleware: JWT, role, object authorization, Redis session, Redis blacklist
  ↓
Handler
  ↓
Service / Business Logic
  ↓
Repository
  ↓
MySQL
```

Redis digunakan sebagai pendukung performa dan keamanan:

```text
Redis Cache
Redis Login Rate Limiter
Redis Token Blacklist
Redis Idle Session Timeout
Redis Lock untuk race condition
Cache Invalidation setelah operasi write
```

## Tech Stack

| Komponen               | Teknologi                  |
| ---------------------- | -------------------------- |
| Language               | Go                         |
| HTTP Framework         | Gin                        |
| Database               | MySQL / MariaDB            |
| Cache & Security Store | Redis / Redis Cloud        |
| Auth                   | JWT Bearer Token           |
| API Testing            | Postman                    |
| Architecture           | Handler-Service-Repository |

## Fitur Utama

- JWT authentication.
- Role-based authorization: `admin`, `trainer`, `talent`.
- Object-level authorization untuk mencegah user mengakses data yang bukan haknya.
- Redis login rate limiter.
- Redis token blacklist saat logout.
- Redis idle session timeout 30 menit.
- Redis response cache untuk endpoint GET tertentu.
- Redis cache invalidation setelah data berubah.
- Redis lock untuk mencegah double submit/resubmit.
- MySQL transaction untuk operasi penting.
- `SELECT ... FOR UPDATE` untuk mengurangi race condition pada review/resubmit.
- Unique constraint `(assignment_id, talent_id)` untuk mencegah duplicate submission.
- Negative test Postman untuk role authorization dan object-level authorization.

## Struktur Project

```text
cmd/
  api/
    main.go

internal/
  cache/
    redis.go

  config/
    config.go

  database/
    mysql.go

  middleware/
    auth_middleware.go
    role_middleware.go
    redis_cache_middleware.go
    login_rate_limiter.go
    jwt_blacklist_middleware.go
    redis_lock_middleware.go
    idle_session_middleware.go
    object_authorization_middleware.go

  handler/
    auth_handler.go
    logout_handler.go
    class_handler.go
    assignment_handler.go
    submission_handler.go
    review_handler.go

  service/
    auth_service.go
    class_service.go
    assignment_service.go
    submission_service.go
    review_service.go

  repository/
    user_repository.go
    class_repository.go
    assignment_repository.go
    submission_repository.go
    review_repository.go
    dbtx.go

  model/
    user.go
    class.go
    assignment.go
    submission.go

  dto/
    auth_dto.go
    class_dto.go
    assignment_dto.go
    submission_dto.go

  utils/
    jwt.go
    response.go
    password.go

db/
  migrations/
    000001_init_schema.up.sql
    000001_init_schema.down.sql

README.md
.env.example
go.mod
```

## Setup Project

### 1. Copy env

```bash
cp .env.example .env
```

Contoh isi `.env`:

```env
APP_PORT=9099
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=bootcamp_assignment
JWT_SECRET=your_secret_key

REDIS_ADDR=your-redis-host:your-redis-port
REDIS_USERNAME=default
REDIS_PASSWORD=your_redis_password
REDIS_DB=0
```

Jika menggunakan Redis Cloud, ambil nilai `REDIS_ADDR`, `REDIS_USERNAME`, dan `REDIS_PASSWORD` dari menu **Connect** pada dashboard Redis Cloud.

### 2. Buat database

```sql
CREATE DATABASE bootcamp_assignment;
```

### 3. Jalankan migration

Linux/macOS:

```bash
mysql -u root -p bootcamp_assignment < db/migrations/000001_init_schema.up.sql
```

PowerShell:

```powershell
Get-Content db/migrations/000001_init_schema.up.sql | mysql -u root -p bootcamp_assignment
```

### 4. Install dependency Go

```bash
go mod tidy
```

Jika Redis client belum ada:

```bash
go get github.com/redis/go-redis/v9
```

Jika JWT v5 belum ada:

```bash
go get github.com/golang-jwt/jwt/v5
```

### 5. Jalankan aplikasi

```bash
go run ./cmd/api
```

Health check:

```http
GET http://localhost:9099/health
```

Expected response:

```json
{
  "status": "ok",
  "message": "Bootcamp Assignment API is running"
}
```

## Akun Demo

Password semua akun demo: `password123`.

| Role      | Email                  | Keterangan            |
| --------- | ---------------------- | --------------------- |
| Admin     | `admin@example.com`    | Akses penuh           |
| Trainer A | `trainer@example.com`  | Trainer untuk Class A |
| Trainer B | `trainer2@example.com` | Trainer untuk Class B |
| Talent A  | `talent@example.com`   | Talent untuk Class A  |
| Talent B  | `talent2@example.com`  | Talent untuk Class B  |

## Seed Data untuk Negative Test

Agar object-level authorization bisa dites dengan jelas, gunakan data seperti ini:

```text
Admin       = id 1
Trainer A   = id 2
Trainer B   = id 3
Talent A    = id 4
Talent B    = id 5

Class A     = id 1 -> Trainer A, Talent A
Class B     = id 2 -> Trainer B, Talent B

Assignment A = id 1 -> Class A
Assignment B = id 2 -> Class B

Submission A = id 1 -> Assignment A, Talent A
Submission B = id 2 -> Assignment B, Talent B
```

Variable Postman yang disarankan:

```text
base_url=http://localhost:9099

class_a_id=1
class_b_id=2

assignment_a_id=1
assignment_b_id=2

submission_a_id=1
submission_b_id=2

talent_a_id=4
talent_b_id=5

trainer_a_id=2
trainer_b_id=3

admin_token=
trainer_a_token=
trainer_b_token=
talent_a_token=
talent_b_token=
```

## Endpoint List

| Area       | Method | Endpoint                                        | Role                                 | Keterangan                         |
| ---------- | ------ | ----------------------------------------------- | ------------------------------------ | ---------------------------------- |
| Health     | GET    | `/health`                                       | Public                               | Cek server hidup                   |
| Auth       | POST   | `/api/v1/auth/login`                            | Public                               | Login dan menghasilkan JWT         |
| Auth       | GET    | `/api/v1/auth/me`                               | Login                                | Mengambil profil user dari token   |
| Auth       | POST   | `/api/v1/auth/logout`                           | Login                                | Logout dan blacklist token         |
| Class      | POST   | `/api/v1/classes`                               | Admin                                | Membuat class/batch                |
| Class      | GET    | `/api/v1/classes?page=1&limit=10`               | Login                                | List class dengan pagination       |
| Class      | GET    | `/api/v1/classes/:id`                           | Login + Object Access                | Detail class                       |
| Class      | POST   | `/api/v1/classes/:id/trainers`                  | Admin                                | Assign trainer ke class            |
| Class      | POST   | `/api/v1/classes/:id/talents`                   | Admin                                | Assign talent ke class             |
| Class      | GET    | `/api/v1/classes/:id/trainers`                  | Admin/Trainer + Object Access        | List trainer dalam class           |
| Class      | GET    | `/api/v1/classes/:id/talents`                   | Admin/Trainer + Object Access        | List talent dalam class            |
| Assignment | POST   | `/api/v1/classes/:id/assignments`               | Admin/Trainer                        | Membuat assignment                 |
| Assignment | GET    | `/api/v1/classes/:id/assignments`               | Login + Object Access                | List assignment pada class         |
| Assignment | GET    | `/api/v1/assignments/:id`                       | Login + Object Access                | Detail assignment                  |
| Assignment | PUT    | `/api/v1/assignments/:id`                       | Admin/Trainer + Object Access        | Update assignment                  |
| Assignment | PATCH  | `/api/v1/assignments/:id/close`                 | Admin/Trainer + Object Access        | Close assignment                   |
| Submission | POST   | `/api/v1/assignments/:assignmentId/submissions` | Talent + Object Access               | Submit assignment                  |
| Submission | PUT    | `/api/v1/submissions/:id`                       | Talent + Object Access               | Resubmit jika `revision_required`  |
| Submission | GET    | `/api/v1/submissions/me`                        | Talent                               | List submission milik talent login |
| Submission | GET    | `/api/v1/classes/:id/submissions`               | Admin/Trainer + Object Access        | List submission class              |
| Review     | POST   | `/api/v1/submissions/:id/review`                | Admin/Trainer + Object Access        | Review dengan score dan feedback   |
| Review     | POST   | `/api/v1/submissions/:id/request-revision`      | Admin/Trainer + Object Access        | Meminta revisi                     |
| Progress   | GET    | `/api/v1/classes/:id/progress`                  | Admin/Trainer + Object Access        | Progress summary per class         |
| Progress   | GET    | `/api/v1/talents/:talentId/progress`            | Admin/Trainer/Talent + Object Access | Progress summary per talent        |

## Alur API dari Awal sampai Akhir

### 1. Server start

```text
Load config dari .env
↓
Connect MySQL
↓
Inisialisasi repository
↓
Inisialisasi service
↓
Inisialisasi handler
↓
Connect Redis
↓
Register route Gin
↓
Server berjalan di APP_PORT
```

### 2. Login

```http
POST /api/v1/auth/login
```

Alur login:

```text
Request login
↓
Redis LoginRateLimiter mengecek batas login
↓
Auth service validasi email/password
↓
JWT dibuat
↓
Redis idle session dibuat selama 30 menit
↓
Token dikirim ke client
```

### 3. Request protected API

Setiap request protected harus membawa header:

```http
Authorization: Bearer <token>
```

Alur middleware protected:

```text
TokenBlacklistMiddleware
↓
AuthMiddleware JWT
↓
IdleSessionMiddleware Redis
↓
Role Authorization
↓
Object-Level Authorization
↓
Handler
↓
Service
↓
Repository
↓
MySQL
```

### 4. Admin flow

```text
Admin login
↓
Create class
↓
Assign trainer ke class
↓
Assign talent ke class
↓
Melihat class, assignment, submission, dan progress
```

### 5. Trainer flow

```text
Trainer login
↓
Melihat class yang dia handle
↓
Membuat assignment di class tersebut
↓
Melihat submission talent
↓
Review submission / request revision
↓
Melihat progress class dan talent
```

### 6. Talent flow

```text
Talent login
↓
Melihat class tempat dia terdaftar
↓
Melihat assignment
↓
Submit assignment
↓
Resubmit jika revision_required
↓
Melihat progress miliknya sendiri
```

### 7. Logout

```text
User hit POST /auth/logout
↓
Token dimasukkan ke Redis blacklist
↓
Idle session dihapus
↓
Token yang sama tidak bisa dipakai lagi
```

## Authorization Design

### Role-Based Authorization

Role digunakan untuk membatasi fitur berdasarkan jenis user:

| Role    | Akses                                                              |
| ------- | ------------------------------------------------------------------ |
| Admin   | Mengelola class, assign trainer/talent, melihat semua data, review |
| Trainer | Mengelola assignment dan review pada class yang dia handle         |
| Talent  | Melihat assignment, submit, resubmit, dan melihat progress sendiri |

### Object-Level Authorization

Object-level authorization digunakan untuk memastikan user tidak hanya punya role yang benar, tapi juga punya relasi dengan data yang diakses.

Contoh aturan:

```text
Trainer hanya boleh akses class yang dia handle.
Trainer hanya boleh review submission dari class yang dia handle.
Talent hanya boleh akses class tempat dia terdaftar.
Talent hanya boleh submit assignment dari class tempat dia terdaftar.
Talent hanya boleh resubmit submission miliknya sendiri.
Talent hanya boleh melihat progress miliknya sendiri.
Admin boleh mengakses semua data.
```

Middleware yang digunakan:

```text
RequireClassAccess
RequireAssignmentAccess
RequireSubmissionAccess
RequireTalentProgressAccess
```

## Redis Usage

Redis digunakan untuk lima kebutuhan utama:

```text
1. Login rate limiter
2. Token blacklist saat logout
3. Idle session timeout
4. Response cache
5. Distributed lock untuk race condition
```

### Daftar Endpoint yang Menggunakan Redis

| Endpoint                                             | Redis Usage                            | TTL / Window                       | Keterangan                                           |
| ---------------------------------------------------- | -------------------------------------- | ---------------------------------- | ---------------------------------------------------- |
| `POST /api/v1/auth/login`                            | Rate limiter + create idle session     | 3 request / 5 menit, idle 30 menit | Membatasi brute force login dan membuat session idle |
| `POST /api/v1/auth/logout`                           | Token blacklist + delete idle session  | Sampai JWT expired                 | Token logout tidak bisa digunakan lagi               |
| Semua protected endpoint                             | Blacklist check + idle session refresh | 30 menit                           | Token ditolak jika logout atau idle expired          |
| `GET /api/v1/classes`                                | Response cache                         | 15 menit                           | Cache list class                                     |
| `GET /api/v1/classes/:id`                            | Response cache                         | 15 menit                           | Cache detail class                                   |
| `GET /api/v1/classes/:id/assignments`                | Response cache                         | 10 menit                           | Cache list assignment per class                      |
| `GET /api/v1/assignments/:id`                        | Response cache                         | 10 menit                           | Cache detail assignment                              |
| `GET /api/v1/classes/:id/progress`                   | Response cache                         | 5 menit                            | Cache progress class                                 |
| `GET /api/v1/talents/:talentId/progress`             | Response cache                         | 5 menit                            | Cache progress talent                                |
| `POST /api/v1/classes`                               | Cache invalidation                     | Setelah write sukses               | Hapus cache class                                    |
| `POST /api/v1/classes/:id/trainers`                  | Cache invalidation                     | Setelah write sukses               | Hapus cache class                                    |
| `POST /api/v1/classes/:id/talents`                   | Cache invalidation                     | Setelah write sukses               | Hapus cache class                                    |
| `POST /api/v1/classes/:id/assignments`               | Cache invalidation                     | Setelah write sukses               | Hapus cache assignment/progress                      |
| `PUT /api/v1/assignments/:id`                        | Cache invalidation                     | Setelah write sukses               | Hapus cache assignment/progress                      |
| `PATCH /api/v1/assignments/:id/close`                | Cache invalidation                     | Setelah write sukses               | Hapus cache assignment/progress                      |
| `POST /api/v1/assignments/:assignmentId/submissions` | Redis lock + invalidation              | Lock 30 detik                      | Mencegah double submit dan refresh progress          |
| `PUT /api/v1/submissions/:id`                        | Redis lock + invalidation              | Lock 30 detik                      | Mencegah double resubmit dan refresh progress        |
| `POST /api/v1/submissions/:id/review`                | Cache invalidation                     | Setelah write sukses               | Refresh submission/progress                          |
| `POST /api/v1/submissions/:id/request-revision`      | Cache invalidation                     | Setelah write sukses               | Refresh submission/progress                          |

### Endpoint yang Tidak Dicache

Endpoint write/action tidak dicache karena berisi perubahan data atau aksi security:

```text
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/classes
POST /api/v1/classes/:id/trainers
POST /api/v1/classes/:id/talents
POST /api/v1/classes/:id/assignments
PUT  /api/v1/assignments/:id
PATCH /api/v1/assignments/:id/close
POST /api/v1/assignments/:assignmentId/submissions
PUT  /api/v1/submissions/:id
POST /api/v1/submissions/:id/review
POST /api/v1/submissions/:id/request-revision
```

Untuk endpoint write, Redis dipakai untuk rate limiter, token/session, lock, atau cache invalidation, bukan response cache.

## Redis Key Pattern

| Fitur            | Pattern Key                  | Contoh                                       |
| ---------------- | ---------------------------- | -------------------------------------------- |
| Response cache   | `cache:*:GET:/api/v1/...`    | `cache:user:1:GET:/api/v1/classes`           |
| Login rate limit | `rate_limit:login:<ip>`      | `rate_limit:login:127.0.0.1`                 |
| Token blacklist  | `jwt_blacklist:<token_hash>` | `jwt_blacklist:a3f...`                       |
| Idle session     | `session_idle:<token_hash>`  | `session_idle:b91...`                        |
| Submission lock  | `lock:submission:*`          | `lock:submission:create:assignment:1:user:4` |

Command cek Redis:

```redis
SCAN 0 MATCH cache:* COUNT 100
SCAN 0 MATCH rate_limit:login:* COUNT 100
SCAN 0 MATCH jwt_blacklist:* COUNT 100
SCAN 0 MATCH session_idle:* COUNT 100
SCAN 0 MATCH lock:submission:* COUNT 100
TTL "nama_key"
GET "nama_key"
```

## Transaction & Concurrency

Transaction digunakan pada:

```text
Submit assignment: insert submission + insert history.
Resubmit assignment: update submission + insert history.
Review submission: lock submission, insert review, update status, insert history.
Request revision: update status menjadi revision_required + insert history.
```

Concurrency protection:

| Mekanisme                                      | Fungsi                                                             |
| ---------------------------------------------- | ------------------------------------------------------------------ |
| Redis Lock                                     | Mencegah double submit/resubmit dalam request bersamaan            |
| MySQL Transaction                              | Menjaga proses write tetap atomic                                  |
| SELECT ... FOR UPDATE                          | Mengunci row saat review/resubmit                                  |
| Unique Constraint `(assignment_id, talent_id)` | Mencegah satu talent submit assignment yang sama lebih dari sekali |

Catatan: backend ini tidak memakai mutex untuk mencegah race condition antar request, karena mutex hanya berlaku pada satu proses aplikasi. Untuk backend yang bisa di-scale ke beberapa instance, Redis lock dan database constraint lebih tepat.

## Postman Test

### Login dan simpan token

Contoh script pada tab **Tests** untuk login Trainer A:

```javascript
const json = pm.response.json();

const token = json?.data?.token || json?.token || json?.access_token;

pm.environment.set("trainer_a_token", token);

pm.test("Login trainer A success", function () {
  pm.response.to.have.status(200);
  pm.expect(token).to.be.a("string");
});
```

### Positive test

| Test                        | Request                                | Token             | Expected |
| --------------------------- | -------------------------------------- | ----------------- | -------- |
| Trainer A akses Class A     | `GET /classes/{{class_a_id}}`          | `trainer_a_token` | 200      |
| Trainer B akses Class B     | `GET /classes/{{class_b_id}}`          | `trainer_b_token` | 200      |
| Talent A akses Assignment A | `GET /assignments/{{assignment_a_id}}` | `talent_a_token`  | 200      |
| Talent B akses Assignment B | `GET /assignments/{{assignment_b_id}}` | `talent_b_token`  | 200      |

### Negative object-level authorization test

| Test                                    | Request                                                  | Token             | Expected |
| --------------------------------------- | -------------------------------------------------------- | ----------------- | -------- |
| Trainer A akses Class B                 | `GET /classes/{{class_b_id}}`                            | `trainer_a_token` | 403      |
| Talent A akses Class B                  | `GET /classes/{{class_b_id}}`                            | `talent_a_token`  | 403      |
| Trainer A lihat assignments Class B     | `GET /classes/{{class_b_id}}/assignments`                | `trainer_a_token` | 403      |
| Talent A lihat Assignment B             | `GET /assignments/{{assignment_b_id}}`                   | `talent_a_token`  | 403      |
| Talent A submit Assignment B            | `POST /assignments/{{assignment_b_id}}/submissions`      | `talent_a_token`  | 403      |
| Talent A resubmit Submission B          | `PUT /submissions/{{submission_b_id}}`                   | `talent_a_token`  | 403      |
| Trainer A review Submission B           | `POST /submissions/{{submission_b_id}}/review`           | `trainer_a_token` | 403      |
| Trainer A request revision Submission B | `POST /submissions/{{submission_b_id}}/request-revision` | `trainer_a_token` | 403      |
| Talent A lihat progress Talent B        | `GET /talents/{{talent_b_id}}/progress`                  | `talent_a_token`  | 403      |

Contoh Postman Tests untuk negative test:

```javascript
pm.test("Should be forbidden", function () {
  pm.response.to.have.status(403);
});

pm.test("Response success false", function () {
  const json = pm.response.json();
  pm.expect(json.success).to.eql(false);
});
```

### Negative auth test

| Test                   | Request                                  | Expected |
| ---------------------- | ---------------------------------------- | -------- |
| Tanpa token            | `GET /api/v1/auth/me`                    | 401      |
| Token sudah logout     | `GET /api/v1/auth/me` memakai token lama | 401      |
| Login lebih dari batas | Login ke-4 dalam 5 menit                 | 429      |

## Security Notes

Backend ini sudah memiliki beberapa proteksi:

```text
JWT Bearer Token
Role Authorization
Object-Level Authorization
Redis Login Rate Limiter
Redis Token Blacklist
Redis Idle Session Timeout
Redis Lock
Database Transaction
Unique Constraint
```

Improvement berikutnya yang bisa ditambahkan:

```text
Request body size limit
Strict CORS
Security headers
Audit log async dengan goroutine worker
Dependency scanning dengan govulncheck
Refresh token rotation
Input validation lebih detail per DTO
Max pagination limit
```

## Catatan Implementasi Penting

Pastikan route berikut sudah sesuai sebelum menjalankan negative test:

### Create assignment harus cek akses class

```go
protected.POST(
    "/classes/:id/assignments",
    middleware.RequireRole("admin", "trainer"),
    middleware.RequireClassAccess(db, "id"),
    middleware.InvalidateCache(...),
    assignmentHandler.CreateAssignment,
)
```

### Resubmit harus cek akses submission, bukan assignment

```go
protected.PUT(
    "/submissions/:id",
    middleware.RequireRole("talent"),
    middleware.RequireSubmissionAccess(db, "id"),
    middleware.RedisLock(redisClient, 30*time.Second, middleware.SubmissionUpdateLockKey),
    middleware.InvalidateCache(...),
    submissionHandler.ResubmitSubmission,
)
```

Kalau masih memakai `RequireAssignmentAccess(db, "id")` pada route `/submissions/:id`, maka `id` akan dianggap sebagai assignment ID, padahal yang dikirim adalah submission ID.

## Interview Explanation

Versi singkat untuk dijelaskan saat interview:

> Backend ini menggunakan Go Gin dengan layered architecture. Request masuk ke router, melewati middleware security seperti JWT authentication, role authorization, Redis blacklist, Redis idle session, dan object-level authorization. Handler bertugas menerima request, service menyimpan business logic, dan repository mengakses MySQL. Redis digunakan untuk rate limiter login, token blacklist saat logout, idle session timeout, response cache, cache invalidation, dan distributed lock untuk mencegah double submit. Untuk data consistency, operasi submit, resubmit, dan review memakai database transaction, row locking, dan unique constraint.
