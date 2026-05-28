# 17 Poin Mapping

Project ini dibuat mengikuti dokumen desain Backend Design Document - Bootcamp Assignment & Submission Platform.

## 1. Problem Summary

Platform untuk mengelola class/batch bootcamp, assignment, submission talent, review trainer, request revision, dan progress tracking.

## 2. Assumptions

- Role utama: admin, trainer, talent.
- Satu class dapat memiliki banyak trainer dan talent.
- Assignment dibuat untuk satu class.
- Talent hanya melihat dan submit assignment dari class tempat dia terdaftar.
- Talent hanya resubmit jika submission berstatus `revision_required`.
- Trainer hanya review submission dari class yang dia handle.
- Admin memiliki akses penuh.
- MySQL InnoDB digunakan untuk foreign key dan transaction.
- Redis dan Docker tidak digunakan pada versi awal.

## 3. Architecture Overview

```text
Client/Postman
  ↓
Gin Handler
  ↓
Service Layer
  ↓
Repository Layer
  ↓
MySQL
```

## 4. API Endpoint List

Semua endpoint pada dokumen sudah ada di project:

- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/classes`
- `GET /api/v1/classes`
- `GET /api/v1/classes/:id`
- `POST /api/v1/classes/:id/trainers`
- `POST /api/v1/classes/:id/talents`
- `POST /api/v1/classes/:classId/assignments`
- `GET /api/v1/classes/:classId/assignments`
- `GET /api/v1/assignments/:id`
- `PUT /api/v1/assignments/:id`
- `PATCH /api/v1/assignments/:id/close`
- `POST /api/v1/assignments/:assignmentId/submissions`
- `PUT /api/v1/submissions/:id`
- `GET /api/v1/submissions/me`
- `GET /api/v1/classes/:classId/submissions`
- `POST /api/v1/submissions/:id/review`
- `POST /api/v1/submissions/:id/request-revision`
- `GET /api/v1/classes/:classId/progress`
- `GET /api/v1/talents/:talentId/progress`

## 5. Database Schema or ERD

Migration membuat tabel:

- `users`
- `classes`
- `class_trainers`
- `class_talents`
- `assignments`
- `submissions`
- `submission_reviews`
- `submission_histories`

## 6. Main Request Flows

### Talent Submit Assignment

Login talent → lihat assignment → submit GitHub/deployment URL → sistem cek membership class → cek deadline → simpan submission → insert history.

### Trainer Review Submission

Login trainer → lihat submission class → review score/feedback → sistem cek trainer handle class → transaction insert review, update status, insert history.

### Request Revision

Trainer request revision → status menjadi `revision_required` → talent bisa resubmit lewat `PUT /api/v1/submissions/:id`.

## 7. Go Project Structure

Struktur sudah dipisah menjadi `handler`, `service`, `repository`, `model`, `dto`, `middleware`, `utils`, dan `db/migrations`.

## 8. Authentication and Authorization Design

Authentication menggunakan JWT. Authorization menggunakan middleware role dan validasi tambahan di service.

## 9. Validation and Error Handling Strategy

Response sukses dan error dibuat konsisten lewat `utils/response.go`.

## 10. Transaction Handling Strategy

Transaction digunakan pada submit, resubmit, review, dan request revision.

## 11. Caching Decision

Belum memakai Redis/cache karena data submission, review, score, feedback, dan revision harus fresh.

## 12. Concurrency and Race Condition Consideration

Review dan resubmit memakai `SELECT ... FOR UPDATE` agar status terbaru dibaca dan row terkunci saat diubah.

## 13. Background Job Consideration

Belum ada background job pada versi awal. Status `late` dihitung saat submit/resubmit dan progress dihitung langsung dari database.

## 14. Testing Plan

Test bisa dilakukan dengan Postman memakai `docs/postman_examples.md`.

## 15. Deployment Consideration

Deployment sederhana: install MySQL, jalankan migration, build binary Go, jalankan via systemd, opsional Nginx reverse proxy.

## 16. Trade-offs and Limitations

Monolith sederhana dipilih agar mudah dijelaskan. Tidak memakai Redis dan Docker untuk menghindari kompleksitas versi awal.

## 17. AI / Tools Used and What I Personally Verified

AI digunakan untuk membantu menyusun struktur dan contoh kode. Yang diverifikasi: role access, endpoint, database relationship, transaction flow, caching decision, dan concurrency handling.
