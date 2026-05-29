# Bootcamp Assignment & Submission Platform

Project backend Go Gin + MySQL untuk assessment **Trainer Candidate Technical Assessment**.
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
- Tanpa Redis dan tanpa Docker untuk versi awal sesuai dokumen.

## Struktur Project

```text
cmd/
  api/
    main.go
internal/
  config/
    config.go
  database/
    mysql.go
  middleware/
    auth_middleware.go
    role_middleware.go
  handler/
    auth_handler.go
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

Isi `.env` sesuai database lokal.

```env
APP_PORT=9099
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=bootcamp_assignment
JWT_SECRET=your_secret_key
```

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

### 4. Jalankan aplikasi

```bash
go run ./cmd/api
```

Health check:

```http
GET http://localhost:9099/health
```

## Akun Demo

Password semua akun demo: `password123`.

| Role | Email |
|---|---|
| Admin | `admin@example.com` |
| Trainer | `trainer@example.com` |
| Talent | `talent@example.com` |

## Endpoint List

| Area | Method | Endpoint | Role | Keterangan |
|---|---|---|---|---|
| Health | GET | `/health` | Public | Cek server hidup |
| Auth | POST | `/api/v1/auth/login` | Public | Login dan menghasilkan JWT |
| Auth | GET | `/api/v1/auth/me` | Login | Mengambil profil user dari token |
| Class | POST | `/api/v1/classes` | Admin | Membuat class/batch |
| Class | GET | `/api/v1/classes?page=1&limit=10` | Login | List class dengan pagination |
| Class | GET | `/api/v1/classes/:id` | Login | Detail class |
| Class | POST | `/api/v1/classes/:id/trainers` | Admin | Assign trainer ke class |
| Class | POST | `/api/v1/classes/:id/talents` | Admin | Assign talent ke class |
| Assignment | POST | `/api/v1/classes/:classId/assignments` | Admin/Trainer | Membuat assignment |
| Assignment | GET | `/api/v1/classes/:classId/assignments` | Login | List assignment pada class |
| Assignment | GET | `/api/v1/assignments/:id` | Login | Detail assignment |
| Assignment | PUT | `/api/v1/assignments/:id` | Admin/Trainer | Update assignment |
| Assignment | PATCH | `/api/v1/assignments/:id/close` | Admin/Trainer | Close assignment |
| Submission | POST | `/api/v1/assignments/:assignmentId/submissions` | Talent | Submit assignment |
| Submission | PUT | `/api/v1/submissions/:id` | Talent | Resubmit jika `revision_required` |
| Submission | GET | `/api/v1/submissions/me` | Talent | List submission milik talent login |
| Submission | GET | `/api/v1/classes/:classId/submissions` | Admin/Trainer | List submission class |
| Review | POST | `/api/v1/submissions/:id/review` | Admin/Trainer | Review dengan score dan feedback |
| Review | POST | `/api/v1/submissions/:id/request-revision` | Admin/Trainer | Meminta revisi |
| Progress | GET | `/api/v1/classes/:classId/progress` | Admin/Trainer | Progress summary per class |
| Progress | GET | `/api/v1/talents/:talentId/progress` | Admin/Trainer/Talent | Progress summary per talent |

## Catatan Role

- Admin memiliki akses penuh untuk class, assignment, review, dan progress.
- Trainer hanya boleh mengelola assignment dan review submission dari class yang dia handle.
- Talent hanya boleh melihat assignment dari class tempat dia terdaftar dan submit/resubmit tugas miliknya sendiri.

## Transaction & Concurrency

Transaction digunakan pada:

- Submit assignment: insert submission + insert history.
- Resubmit assignment: update submission + insert history.
- Review submission: lock submission, insert review, update status, insert history.
- Request revision: sama seperti review, tetapi status dipaksa menjadi `revision_required`.

Pada review dan resubmit, sistem memakai `SELECT ... FOR UPDATE` agar dua proses tidak mengubah submission yang sama secara bersamaan.

## Postman

Lihat contoh lengkap di:

```text
docs/postman_examples.md
```
