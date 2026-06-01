# Bootcamp Assignment & Submission Platform

Project backend **Go Gin + MySQL + Redis** untuk assessment **Trainer Candidate Technical Assessment**.
Project ini dibuat agar sesuai dengan dokumen desain: `Backend Design Document - Bootcamp Assignment & Submission Platform`.

## Fitur Utama

- JWT authentication.
- Role authorization: `admin`, `trainer`, `talent`.
- Class/batch management.
- Assign trainer dan talent ke class.
- Assignment management: create, list, detail, update, close.
- Submission management: submit, resubmit saat `revision_required`, list submission.
- Review management: score, feedback, request revision.
- Progress tracking per class dan per talent.
- MySQL transaction untuk operasi penting.
- `SELECT ... FOR UPDATE` untuk mengurangi race condition pada review/resubmit.
- Redis Cloud / Redis-compatible server untuk caching, rate limiter, token blacklist, dan distributed lock sederhana.
- Response cache untuk endpoint read-heavy seperti class list, class detail, assignment list, assignment detail, dan progress tracking.
- Cache invalidation setelah operasi write seperti create class, assign trainer/talent, create/update/close assignment, submit/resubmit, review, dan request revision.
- Login rate limiter untuk mengurangi risiko brute force.
- Logout dengan JWT blacklist di Redis agar token yang sudah logout tidak bisa dipakai kembali.
- Redis lock untuk mencegah double submit/resubmit dari talent dalam waktu bersamaan.

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
    login_rate_limiter.go
    redis_cache_middleware.go
    jwt_blacklist_middleware.go
    redis_lock_middleware.go
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
```

## Setup

### 1. Copy env

```bash
cp .env.example .env
```

Isi `.env` sesuai database lokal dan Redis yang digunakan.

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

Catatan:

- Untuk Redis Cloud, gunakan host, port, username, dan password dari dashboard Redis Cloud.
- Jangan upload file `.env` ke GitHub karena berisi credential database dan Redis.

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

Jika package Redis belum terinstall, jalankan:

```bash
go get github.com/redis/go-redis/v9
```

Jika menggunakan JWT v5 untuk membaca expiry token blacklist, jalankan:

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

Jika Redis berhasil terkoneksi, log aplikasi akan menampilkan pesan seperti:

```text
connected to Redis successfully
```

## Akun Demo

Password semua akun demo: `password123`.

| Role    | Email                 |
| ------- | --------------------- |
| Admin   | `admin@example.com`   |
| Trainer | `trainer@example.com` |
| Talent  | `talent@example.com`  |

## Endpoint List

| Area       | Method | Endpoint                                        | Role                 | Keterangan                                          |
| ---------- | ------ | ----------------------------------------------- | -------------------- | --------------------------------------------------- |
| Health     | GET    | `/health`                                       | Public               | Cek server hidup                                    |
| Auth       | POST   | `/api/v1/auth/login`                            | Public               | Login dan menghasilkan JWT                          |
| Auth       | GET    | `/api/v1/auth/me`                               | Login                | Mengambil profil user dari token                    |
| Auth       | POST   | `/api/v1/auth/logout`                           | Login                | Logout dan memasukkan JWT ke Redis blacklist        |
| Class      | POST   | `/api/v1/classes`                               | Admin                | Membuat class/batch                                 |
| Class      | GET    | `/api/v1/classes?page=1&limit=10`               | Login                | List class dengan pagination dan cache              |
| Class      | GET    | `/api/v1/classes/:id`                           | Login                | Detail class dengan cache                           |
| Class      | POST   | `/api/v1/classes/:id/trainers`                  | Admin                | Assign trainer ke class                             |
| Class      | POST   | `/api/v1/classes/:id/talents`                   | Admin                | Assign talent ke class                              |
| Class      | GET    | `/api/v1/classes/:id/trainers`                  | Admin/Trainer        | List trainer pada class                             |
| Class      | GET    | `/api/v1/classes/:id/talents`                   | Admin/Trainer        | List talent pada class                              |
| Assignment | POST   | `/api/v1/classes/:id/assignments`               | Admin/Trainer        | Membuat assignment                                  |
| Assignment | GET    | `/api/v1/classes/:id/assignments`               | Login                | List assignment pada class dengan cache             |
| Assignment | GET    | `/api/v1/assignments/:id`                       | Login                | Detail assignment dengan cache                      |
| Assignment | PUT    | `/api/v1/assignments/:id`                       | Admin/Trainer        | Update assignment                                   |
| Assignment | PATCH  | `/api/v1/assignments/:id/close`                 | Admin/Trainer        | Close assignment                                    |
| Submission | POST   | `/api/v1/assignments/:assignmentId/submissions` | Talent               | Submit assignment dengan Redis lock                 |
| Submission | PUT    | `/api/v1/submissions/:id`                       | Talent               | Resubmit dengan Redis lock jika `revision_required` |
| Submission | GET    | `/api/v1/submissions/me`                        | Talent               | List submission milik talent login                  |
| Submission | GET    | `/api/v1/classes/:id/submissions`               | Admin/Trainer        | List submission class                               |
| Review     | POST   | `/api/v1/submissions/:id/review`                | Admin/Trainer        | Review dengan score dan feedback                    |
| Review     | POST   | `/api/v1/submissions/:id/request-revision`      | Admin/Trainer        | Meminta revisi                                      |
| Progress   | GET    | `/api/v1/classes/:id/progress`                  | Admin/Trainer        | Progress summary per class dengan cache             |
| Progress   | GET    | `/api/v1/talents/:talentId/progress`            | Admin/Trainer/Talent | Progress summary per talent dengan cache            |

## Catatan Role

- Admin memiliki akses penuh untuk class, assignment, review, dan progress.
- Trainer hanya boleh mengelola assignment dan review submission dari class yang dia handle.
- Talent hanya boleh melihat assignment dari class tempat dia terdaftar dan submit/resubmit tugas miliknya sendiri.

## Redis Usage

Redis digunakan untuk lima kebutuhan utama:

### 1. Response Cache

Endpoint GET yang sering dibaca disimpan sementara di Redis agar tidak selalu query ke MySQL.

| Endpoint                                 | TTL      |
| ---------------------------------------- | -------- |
| `GET /api/v1/classes`                    | 15 menit |
| `GET /api/v1/classes/:id`                | 15 menit |
| `GET /api/v1/classes/:id/assignments`    | 10 menit |
| `GET /api/v1/assignments/:id`            | 10 menit |
| `GET /api/v1/classes/:id/progress`       | 5 menit  |
| `GET /api/v1/talents/:talentId/progress` | 5 menit  |

Response cache mengembalikan header:

```text
X-Cache: MISS
```

jika data belum ada di Redis dan harus diambil dari MySQL.

```text
X-Cache: HIT
```

jika data sudah ada di Redis dan langsung dikembalikan dari cache.

### 2. Cache Invalidation

Setiap operasi write akan menghapus cache terkait agar data yang dibaca user tetap fresh.

Contoh operasi yang melakukan invalidation:

- Create class.
- Assign trainer ke class.
- Assign talent ke class.
- Create assignment.
- Update assignment.
- Close assignment.
- Submit assignment.
- Resubmit assignment.
- Review submission.
- Request revision.

### 3. Login Rate Limiter

Endpoint login dibatasi menggunakan Redis:

```text
POST /api/v1/auth/login
```

Konfigurasi saat ini:

```text
3 request per 1 menit per IP
```

Jika limit terlampaui, API akan mengembalikan response:

```http
429 Too Many Requests
```

### 4. JWT Blacklist Saat Logout

Saat user logout, JWT yang sedang digunakan akan dimasukkan ke Redis blacklist sampai masa expired token habis.

Alurnya:

```text
POST /api/v1/auth/logout
↓
Token JWT di-hash
↓
Hash token disimpan ke Redis dengan TTL sesuai sisa umur JWT
↓
Token yang sama tidak bisa dipakai lagi
```

### 5. Redis Lock untuk Race Condition

Redis lock digunakan pada submit dan resubmit agar request yang sama tidak diproses bersamaan.

Endpoint yang memakai Redis lock:

```text
POST /api/v1/assignments/:assignmentId/submissions
PUT  /api/v1/submissions/:id
```

Jika request yang sama sedang diproses, API mengembalikan:

```http
409 Conflict
```

## Daftar Endpoint Berdasarkan Penggunaan Redis

Catatan: semua endpoint di dalam protected route tetap melewati `TokenBlacklistMiddleware`, sehingga token yang sudah logout akan dicek ke Redis terlebih dahulu. Tabel di bawah menjelaskan penggunaan Redis khusus seperti cache, rate limiter, lock, dan cache invalidation.

### Endpoint yang Menggunakan Redis

| Method | Endpoint                                        | Redis yang Dipakai              | Tujuan                                                                            |
| ------ | ----------------------------------------------- | ------------------------------- | --------------------------------------------------------------------------------- |
| POST   | `/api/v1/auth/login`                            | `LoginRateLimiter`              | Membatasi percobaan login, saat ini 3 request per 1 menit per IP.                 |
| POST   | `/api/v1/auth/logout`                           | JWT blacklist                   | Menyimpan hash JWT ke Redis agar token yang sudah logout tidak bisa dipakai lagi. |
| GET    | `/api/v1/classes`                               | Response cache                  | Menyimpan list class selama 15 menit.                                             |
| GET    | `/api/v1/classes/:id`                           | Response cache                  | Menyimpan detail class selama 15 menit.                                           |
| POST   | `/api/v1/classes`                               | Cache invalidation              | Menghapus cache class setelah class baru dibuat.                                  |
| POST   | `/api/v1/classes/:id/trainers`                  | Cache invalidation              | Menghapus cache class terkait setelah trainer ditambahkan.                        |
| POST   | `/api/v1/classes/:id/talents`                   | Cache invalidation              | Menghapus cache class terkait setelah talent ditambahkan.                         |
| GET    | `/api/v1/classes/:id/assignments`               | Response cache                  | Menyimpan list assignment class selama 10 menit.                                  |
| GET    | `/api/v1/assignments/:id`                       | Response cache                  | Menyimpan detail assignment selama 10 menit.                                      |
| POST   | `/api/v1/classes/:id/assignments`               | Cache invalidation              | Menghapus cache assignment dan progress setelah assignment dibuat.                |
| PUT    | `/api/v1/assignments/:id`                       | Cache invalidation              | Menghapus cache assignment dan progress setelah assignment diupdate.              |
| PATCH  | `/api/v1/assignments/:id/close`                 | Cache invalidation              | Menghapus cache assignment dan progress setelah assignment ditutup.               |
| POST   | `/api/v1/assignments/:assignmentId/submissions` | Redis lock + cache invalidation | Mencegah double submit dan menghapus cache progress/submission terkait.           |
| PUT    | `/api/v1/submissions/:id`                       | Redis lock + cache invalidation | Mencegah double resubmit dan menghapus cache progress/submission terkait.         |
| POST   | `/api/v1/submissions/:id/review`                | Cache invalidation              | Menghapus cache progress/submission setelah review diberikan.                     |
| POST   | `/api/v1/submissions/:id/request-revision`      | Cache invalidation              | Menghapus cache progress/submission setelah request revision dibuat.              |
| GET    | `/api/v1/classes/:id/progress`                  | Response cache                  | Menyimpan progress class selama 5 menit.                                          |
| GET    | `/api/v1/talents/:talentId/progress`            | Response cache                  | Menyimpan progress talent selama 5 menit.                                         |

### Endpoint yang Tidak Diberi Redis Khusus

Endpoint berikut tidak menggunakan response cache, rate limiter, Redis lock, atau cache invalidation khusus. Jika endpoint berada di protected route, endpoint tetap melewati pengecekan JWT blacklist secara global.

| Method | Endpoint                          | Alasan Tidak Diberi Redis Khusus                                                                                                  |
| ------ | --------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| GET    | `/health`                         | Public health check ringan dan tidak membutuhkan cache.                                                                           |
| GET    | `/api/v1/auth/me`                 | Data user dari token/login context, tidak perlu response cache.                                                                   |
| GET    | `/api/v1/classes/:id/trainers`    | Data membership bisa berubah dan biasanya tidak seberat list/progress. Bisa ditambah cache jika traffic tinggi.                   |
| GET    | `/api/v1/classes/:id/talents`     | Data membership bisa berubah dan biasanya tidak seberat list/progress. Bisa ditambah cache jika traffic tinggi.                   |
| GET    | `/api/v1/submissions/me`          | Data pribadi talent dan sering berubah setelah submit/resubmit/review. Untuk versi awal tidak dicache.                            |
| GET    | `/api/v1/classes/:id/submissions` | Data submission class sering berubah saat talent submit dan trainer review. Untuk versi awal tidak dicache agar data lebih fresh. |

### Prinsip Pemakaian Redis

- Endpoint `GET` yang sering dibaca dan relatif aman disimpan sementara diberi response cache.
- Endpoint `POST`, `PUT`, dan `PATCH` tidak dicache, tetapi dipakai untuk menghapus cache yang terdampak.
- Endpoint yang rawan brute force diberi rate limiter.
- Endpoint submit/resubmit diberi Redis lock karena rawan double request.
- Logout memakai token blacklist agar JWT yang belum expired tetap bisa dinonaktifkan.

## Transaction & Concurrency

Transaction digunakan pada:

- Submit assignment: insert submission + insert history.
- Resubmit assignment: update submission + insert history.
- Review submission: lock submission, insert review, update status, insert history.
- Request revision: sama seperti review, tetapi status dipaksa menjadi `revision_required`.

Pada review dan resubmit, sistem memakai `SELECT ... FOR UPDATE` agar dua proses tidak mengubah submission yang sama secara bersamaan.

Selain itu, Redis lock digunakan untuk mencegah request submit/resubmit yang sama diproses secara paralel dari sisi aplikasi.

Untuk proteksi maksimal, tetap disarankan menambahkan unique constraint pada database, misalnya satu talent hanya boleh memiliki satu submission aktif untuk satu assignment:

```sql
ALTER TABLE submissions
ADD UNIQUE KEY uq_assignment_talent (assignment_id, talent_id);
```

## Cara Test Redis

### Test response cache

1. Login dan ambil JWT.
2. Hit endpoint berikut:

```http
GET /api/v1/classes
Authorization: Bearer <token>
```

3. Request pertama akan menghasilkan response header:

```text
X-Cache: MISS
```

4. Request kedua ke endpoint yang sama akan menghasilkan:

```text
X-Cache: HIT
```

### Test rate limiter login

Kirim request login lebih dari 3 kali dalam 1 menit dari IP yang sama:

```http
POST /api/v1/auth/login
```

Request berikutnya akan mendapat:

```http
429 Too Many Requests
```

### Test logout blacklist

1. Login dan ambil JWT.
2. Logout:

```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

3. Pakai token yang sama untuk endpoint protected:

```http
GET /api/v1/auth/me
Authorization: Bearer <token_yang_sudah_logout>
```

4. API akan menolak token tersebut.

### Cek key Redis

Di Redis CLI / Redis Insight:

```redis
SCAN 0 MATCH cache:* COUNT 100
SCAN 0 MATCH rate_limit:login:* COUNT 100
SCAN 0 MATCH jwt_blacklist:* COUNT 100
SCAN 0 MATCH lock:submission:* COUNT 100
```

Cek isi cache:

```redis
GET "nama_key_cache"
```

Cek sisa umur key:

```redis
TTL "nama_key"
```

## Catatan Keamanan

- JWT secret wajib kuat dan tidak boleh di-hardcode.
- Credential Redis dan database wajib disimpan di `.env`.
- File `.env` tidak boleh di-commit ke repository.
- Redis untuk token blacklist menyimpan hash token, bukan token mentah.
- Redis lock membantu mencegah double request, tetapi database constraint tetap diperlukan sebagai lapisan terakhir.
- Rate limiter login membantu mengurangi brute force, tetapi di production bisa ditambah kombinasi IP + email agar lebih akurat.

## Catatan Cache

Cache hanya digunakan untuk endpoint GET yang read-heavy.
Endpoint write seperti POST, PUT, PATCH tidak di-cache, tetapi digunakan untuk menghapus cache yang terdampak.

Jika data diubah langsung dari MySQL tanpa melalui API, Redis tidak otomatis tahu perubahan tersebut. Cache akan hilang setelah TTL habis atau setelah ada endpoint write yang melakukan invalidation.
