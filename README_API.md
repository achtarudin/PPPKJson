# PPPK JSON Exam API

API untuk sistem ujian PPPK dengan fitur random soal per kategori.

## Overview

Sistem ini memungkinkan user mengerjakan ujian PPPK dengan ketentuan:
- Setiap user mendapat soal yang berbeda
- 4 kategori soal: **MANAJERIAL**, **SOSIAL_KULTURAL**, **TEKNIS**, **WAWANCARA**  
- Masing-masing kategori terdiri dari 5 soal (total 20 soal)
- User ID hardcoded dari URL path
- Sistem scoring 1-4 untuk setiap opsi jawaban

## Database Schema

### Tables
1. **exam_sessions** - Sesi ujian per user
2. **exam_questions** - Soal-soal yang di-assign ke sesi ujian tertentu
3. **user_answers** - Jawaban yang diberikan user
4. **exam_results** - Hasil per kategori
5. **exam_summaries** - Ringkasan hasil keseluruhan

## API Endpoints

### Base URL
```
http://localhost:8080
```

### 1. Create/Get Exam Session
```http
GET /exam/{userID}
```

**Description**: Membuat sesi ujian baru atau mengambil sesi yang sedang berlangsung untuk user tertentu.

**Parameters**:
- `userID` (path): ID user (hardcoded, misal: 1234)

**Response**:
```json
{
  "success": true,
  "message": "Exam session ready",
  "data": {
    "session_id": 1,
    "user_id": "1234",
    "session_code": "EXAM_1234_1643356800",
    "status": "NOT_STARTED",
    "expires_at": "2026-01-28T12:00:00Z",
    "duration": 120,
    "questions": [
      {
        "exam_question_id": 1,
        "question_id": 15,
        "category": "MANAJERIAL",
        "order_number": 1,
        "question_text": "Atasan Anda melakukan rekayasa laporan...",
        "options": [
          {
            "id": 59,
            "option_text": "Dalam hati tidak menyetujui hal tersebut"
          },
          {
            "id": 60,
            "option_text": "Hal tersebut sering terjadi di kantor manapun"
          }
        ]
      }
    ],
    "category_stats": [
      {
        "category": "MANAJERIAL",
        "total_questions": 5,
        "answered_count": 0
      }
    ]
  }
}
```

### 2. Start Exam
```http
POST /exam/{userID}/start
```

**Description**: Memulai ujian (mengubah status menjadi IN_PROGRESS).

**Response**:
```json
{
  "success": true,
  "message": "Exam started successfully",
  "data": {
    "session_id": 1,
    "status": "IN_PROGRESS",
    "started_at": "2026-01-28T10:00:00Z"
  }
}
```

### 3. Submit Answer
```http
POST /exam/{userID}/answer
```

**Description**: Mengirim jawaban untuk soal tertentu.

**Request Body**:
```json
{
  "exam_question_id": 1,
  "question_option_id": 59
}
```

**Response**:
```json
{
  "success": true,
  "message": "Answer submitted successfully"
}
```

### 4. Complete Exam
```http
POST /exam/{userID}/complete
```

**Description**: Menyelesaikan ujian dan menghitung hasil.

**Response**:
```json
{
  "success": true,
  "message": "Exam completed successfully",
  "data": {
    "session_id": 1,
    "status": "COMPLETED",
    "completed_at": "2026-01-28T11:30:00Z"
  }
}
```

### 5. Get Results
```http
GET /exam/{userID}/results
```

**Description**: Mengambil hasil ujian user.

**Response**:
```json
{
  "success": true,
  "message": "Exam results retrieved",
  "data": {
    "summary": {
      "exam_session_id": 1,
      "user_id": "1234",
      "total_questions": 20,
      "total_answered": 20,
      "total_score": 65,
      "max_score": 80,
      "overall_percentage": 81.25,
      "overall_grade": "B",
      "is_passed": true,
      "completed_at": "2026-01-28T11:30:00Z"
    },
    "results_by_category": [
      {
        "category": "MANAJERIAL",
        "total_questions": 5,
        "total_answered": 5,
        "total_score": 16,
        "max_score": 20,
        "percentage": 80.0,
        "grade": "B",
        "is_passed": true
      }
    ]
  }
}
```

### 6. Health Check
```http
GET /health
```

**Response**:
```json
{
  "status": "ok",
  "service": "PPPKJson Exam API"
}
```

## Usage Flow

1. **Create Exam Session**: `GET /exam/1234`
   - Sistem akan create sesi baru dengan 20 soal random (5 per kategori)

2. **Start Exam**: `POST /exam/1234/start` 
   - Mulai timer ujian

3. **Submit Answers**: `POST /exam/1234/answer`
   - Kirim jawaban satu per satu

4. **Complete Exam**: `POST /exam/1234/complete`
   - Selesaikan ujian dan hitung hasil

5. **Get Results**: `GET /exam/1234/results`
   - Lihat hasil ujian

## Scoring System

- **Per Opsi**: Score 1-4 (dari JSON data)
- **Per Kategori**: Max score 20 (5 soal × 4 poin max)
- **Total Ujian**: Max score 80 (4 kategori × 20 poin max)
- **Passing Grade**: ≥ 60%
- **Grade Scale**: A (≥90%), B (≥80%), C (≥70%), D (≥60%), E (<60%)

## Development

```bash
# Start database
docker-compose -f compose.dev.yaml up postgres

# Run migrations
make migrate-up

# Seed questions
make db-seed

# Run server
make dev-go-server
# atau
go run ./cmd/server
```

Server akan berjalan di `http://localhost:8080`

## Features

✅ **Random Question Assignment** - Setiap user dapat soal yang berbeda  
✅ **Category-based Exam** - 4 kategori dengan 5 soal masing-masing  
✅ **Session Management** - Tracking status ujian per user  
✅ **Automatic Scoring** - Kalkulasi otomatis hasil per kategori dan keseluruhan  
✅ **Exam Timer** - Durasi ujian dengan auto-expire  
✅ **RESTful API** - Clean API dengan proper HTTP methods  
✅ **Transaction Safety** - Database transactions untuk data consistency