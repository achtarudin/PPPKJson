# Question Score Management - Implementation Guide

## Overview

Fitur Question Score Management memungkinkan administrator untuk melihat dan mengedit score pada table questions dan question_options dengan filter berdasarkan kategori.

## Backend API Endpoints

### 1. Get Categories
```
GET /api/v1/questions/categories
```
**Response:**
```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": ["TEKNIS", "MANAJERIAL", "SOSIAL KULTURAL", "WAWANCARA"]
}
```

### 2. Get Questions by Category
```
GET /api/v1/questions?category={category}
```
**Parameters:**
- `category` (optional): Filter by category name

**Response:**
```json
{
  "success": true,
  "message": "Questions retrieved successfully",
  "data": [
    {
      "id": 136,
      "category": "TEKNIS",
      "question_text": "Tugas utama seorang penata layanan operasional adalah ....",
      "options": [
        {
          "id": 601,
          "question_id": 136,
          "option_text": "Menyusun rencana keuangan perusahaan",
          "score": 4,
          "created_at": "2026-02-02T11:48:02.863783+07:00",
          "updated_at": "2026-02-03T04:52:04.229768566+07:00"
        }
      ],
      "created_at": "2026-02-02T11:48:02.863783+07:00",
      "updated_at": "2026-02-02T11:48:02.863783+07:00"
    }
  ]
}
```

### 3. Update Question Option Score
```
PUT /api/v1/questions/{questionID}/option/{optionID}/score
```
**Parameters:**
- `questionID`: Question ID
- `optionID`: Option ID

**Request Body:**
```json
{
  "score": 4
}
```

**Response:**
```json
{
  "success": true,
  "message": "Score updated successfully",
  "data": {
    "id": 601,
    "question_id": 136,
    "option_text": "Menyusun rencana keuangan perusahaan",
    "score": 4,
    "created_at": "2026-02-02T11:48:02.863783+07:00",
    "updated_at": "2026-02-03T04:52:04.229768566+07:00"
  }
}
```

## Frontend UI Features

### Question Manager Page (`/questions`)

1. **Category Filter**
   - Dropdown dengan semua kategori available
   - Option "All Categories" untuk menampilkan semua questions
   - Refresh button untuk reload data

2. **Questions Display**
   - Card layout dengan responsive design
   - Menampilkan Question ID dan Category badge
   - Question text ditampilkan dengan jelas
   - Options ditampilkan dalam grid 2 kolom

3. **Score Management**
   - Setiap option memiliki dropdown untuk mengubah score (1-4)
   - Score update dilakukan secara real-time via API
   - Visual feedback dengan color coding:
     - Score 4: Green (Success)
     - Score 3: Blue (Primary)  
     - Score 2: Yellow (Warning)
     - Score 1: Red (Danger)

4. **Summary Statistics**
   - Total Questions
   - Total Options
   - Current Category
   - Total Categories

## Navigation

Fitur ini dapat diakses melalui:
1. **Main Navigation**: "Question Manager" menu di navbar
2. **Direct URL**: `http://localhost:5173/questions`

## Technical Implementation

### Backend Structure
- **Handler**: `gin_question_handler.go`
- **DTO**: Added `QuestionManagementResponse` and `QuestionOptionManagementResponse`
- **Mapper**: Convert functions dari model GORM ke DTO
- **Routes**: Registered di `main.go`

### Frontend Structure  
- **Component**: `QuestionManager.jsx`
- **API Service**: Extended `api.js` with `questionAPI`
- **Routes**: Added to `App.jsx`
- **Navigation**: Updated `Layout.jsx`

### Features
- ✅ Category filtering
- ✅ Real-time score updates
- ✅ Responsive design with Bootstrap
- ✅ Error handling and success notifications
- ✅ Loading states and UX feedback
- ✅ RESTful API design
- ✅ Swagger documentation

## Usage

1. **Access Question Manager**: Navigate to `/questions`
2. **Filter by Category**: Select category dari dropdown
3. **Update Scores**: Click dropdown pada option dan pilih score baru (1-4)
4. **Real-time Updates**: Perubahan langsung tersimpan ke database
5. **Visual Feedback**: Alert notifications dan color coding untuk status

## Architecture Benefits

- **Hexagonal Architecture**: Clean separation of concerns
- **DTO Pattern**: Prevents GORM model exposure in API
- **RESTful Design**: Standard HTTP methods dan status codes
- **React Best Practices**: Functional components dengan hooks
- **Bootstrap UI**: Consistent dan responsive design
- **Error Handling**: Comprehensive error responses