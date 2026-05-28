# Postman Examples - Bootcamp Assignment API

Base URL:

```http
http://localhost:9099
```

Jika port di `.env` berbeda, sesuaikan `base_url`.

## Environment Variables

```text
base_url=http://localhost:9099
admin_token=
trainer_token=
talent_token=
class_id=1
assignment_id=1
submission_id=1
talent_id=3
```

## 1. Health Check

```http
GET {{base_url}}/health
```

## 2. Login Admin

```http
POST {{base_url}}/api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

Simpan `data.token` ke `admin_token`.

## 3. Login Trainer

```http
POST {{base_url}}/api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "trainer@example.com",
  "password": "password123"
}
```

Simpan `data.token` ke `trainer_token`.

## 4. Login Talent

```http
POST {{base_url}}/api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "talent@example.com",
  "password": "password123"
}
```

Simpan `data.token` ke `talent_token`.

## 5. Current User

```http
GET {{base_url}}/api/v1/auth/me
Authorization: Bearer {{talent_token}}
```

## 6. Create Class - Admin

```http
POST {{base_url}}/api/v1/classes
Authorization: Bearer {{admin_token}}
Content-Type: application/json
```

```json
{
  "name": "Golang Backend Batch 2",
  "description": "Kelas backend Go Gin dan MySQL",
  "start_date": "2026-06-01",
  "end_date": "2026-07-30"
}
```

## 7. List Classes with Pagination

```http
GET {{base_url}}/api/v1/classes?page=1&limit=10
Authorization: Bearer {{admin_token}}
```

## 8. Class Detail

```http
GET {{base_url}}/api/v1/classes/{{class_id}}
Authorization: Bearer {{admin_token}}
```

## 9. Assign Trainer to Class - Admin

```http
POST {{base_url}}/api/v1/classes/{{class_id}}/trainers
Authorization: Bearer {{admin_token}}
Content-Type: application/json
```

```json
{
  "user_id": 2
}
```

## 10. Assign Talent to Class - Admin

```http
POST {{base_url}}/api/v1/classes/{{class_id}}/talents
Authorization: Bearer {{admin_token}}
Content-Type: application/json
```

```json
{
  "user_id": 3
}
```

## 11. Create Assignment - Admin/Trainer

```http
POST {{base_url}}/api/v1/classes/{{class_id}}/assignments
Authorization: Bearer {{trainer_token}}
Content-Type: application/json
```

```json
{
  "title": "Build REST API CRUD Product",
  "description": "Buat REST API Product dengan Gin dan MySQL",
  "deadline": "2026-05-30T23:59:00+07:00"
}
```

## 12. List Assignments by Class

```http
GET {{base_url}}/api/v1/classes/{{class_id}}/assignments
Authorization: Bearer {{talent_token}}
```

## 13. Assignment Detail

```http
GET {{base_url}}/api/v1/assignments/{{assignment_id}}
Authorization: Bearer {{talent_token}}
```

## 14. Update Assignment - Admin/Trainer

```http
PUT {{base_url}}/api/v1/assignments/{{assignment_id}}
Authorization: Bearer {{trainer_token}}
Content-Type: application/json
```

```json
{
  "title": "Build REST API CRUD Product Updated",
  "description": "Tambahkan validasi request dan response error",
  "deadline": "2026-06-01T23:59:00+07:00",
  "status": "active"
}
```

## 15. Close Assignment - Admin/Trainer

```http
PATCH {{base_url}}/api/v1/assignments/{{assignment_id}}/close
Authorization: Bearer {{trainer_token}}
```

## 16. Submit Assignment - Talent

```http
POST {{base_url}}/api/v1/assignments/{{assignment_id}}/submissions
Authorization: Bearer {{talent_token}}
Content-Type: application/json
```

```json
{
  "github_url": "https://github.com/example/product-api",
  "deployment_url": "https://product-api.example.com",
  "notes": "Saya submit tugas pertama"
}
```

## 17. Resubmit Assignment - Talent

Endpoint ini hanya berhasil jika status submission adalah `revision_required`.

```http
PUT {{base_url}}/api/v1/submissions/{{submission_id}}
Authorization: Bearer {{talent_token}}
Content-Type: application/json
```

```json
{
  "github_url": "https://github.com/example/product-api-fixed",
  "deployment_url": "https://product-api-fixed.example.com",
  "notes": "Saya sudah memperbaiki revisi"
}
```

## 18. My Submissions - Talent

```http
GET {{base_url}}/api/v1/submissions/me
Authorization: Bearer {{talent_token}}
```

## 19. Class Submissions - Admin/Trainer

```http
GET {{base_url}}/api/v1/classes/{{class_id}}/submissions
Authorization: Bearer {{trainer_token}}
```

## 20. Review Submission - Admin/Trainer

```http
POST {{base_url}}/api/v1/submissions/{{submission_id}}/review
Authorization: Bearer {{trainer_token}}
Content-Type: application/json
```

```json
{
  "score": 85,
  "feedback": "API sudah berjalan, tapi validasi input perlu diperbaiki.",
  "status": "reviewed"
}
```

## 21. Request Revision - Admin/Trainer

```http
POST {{base_url}}/api/v1/submissions/{{submission_id}}/request-revision
Authorization: Bearer {{trainer_token}}
Content-Type: application/json
```

```json
{
  "score": 60,
  "feedback": "Tolong tambahkan validasi input dan perbaiki response error.",
  "status": "revision_required"
}
```

## 22. Class Progress - Admin/Trainer

```http
GET {{base_url}}/api/v1/classes/{{class_id}}/progress
Authorization: Bearer {{trainer_token}}
```

## 23. Talent Progress - Admin/Trainer/Talent

```http
GET {{base_url}}/api/v1/talents/{{talent_id}}/progress
Authorization: Bearer {{admin_token}}
```

Talent juga bisa akses endpoint ini, tetapi hanya untuk `talent_id` miliknya sendiri.
