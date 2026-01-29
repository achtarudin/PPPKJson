# PPPK Exam Dashboard Frontend

React frontend untuk sistem ujian PPPK dengan dashboard admin.

## Features

- **Login Page**: Input User ID untuk memulai ujian
- **Exam Board**: Interface ujian dengan 4 pertanyaan (1 per kategori)
- **Results Page**: Hasil ujian dengan detail per kategori  
- **Admin Dashboard**: Daftar semua user dengan status ujian
- **User Detail Modal**: Detail lengkap user dari endpoint dashboard

## API Endpoints Consumed

1. `GET /api/v1/dashboard/users` - Daftar semua user
2. `GET /api/v1/exam/{userID}/dashboard` - Detail dashboard user
3. `GET /api/v1/exam/{userID}` - Buat/ambil session ujian
4. `POST /api/v1/exam/{userID}/start` - Mulai ujian
5. `POST /api/v1/exam/{userID}/answer` - Submit jawaban
6. `POST /api/v1/exam/{userID}/complete` - Selesaikan ujian
7. `GET /api/v1/exam/{userID}/results` - Hasil ujian

## Installation & Setup

```bash
# Install dependencies
npm install

# Install bootstrap-icons if not already installed
npm install bootstrap-icons

# Run development server
npm run dev
```

## Usage

1. **Start Backend Server**: Pastikan Go backend berjalan di `http://localhost:8080`

2. **Start Frontend**: 
   ```bash
   npm run dev
   ```
   Frontend akan berjalan di `http://localhost:5173`

3. **Access Pages**:
   - Home/Login: `http://localhost:5173/`
   - Admin Dashboard: `http://localhost:5173/admin`
   - Exam: `http://localhost:5173/exam/{userID}`
   - Results: `http://localhost:5173/results/{userID}`

## Admin Dashboard Features

- **User List Table**: Menampilkan semua user yang pernah ujian
- **Status Badges**: Color-coded status (COMPLETED, IN_PROGRESS, EXPIRED, NOT_STARTED)
- **Search & Filter**: By User ID dan Status
- **Auto-refresh**: Setiap 30 detik
- **Detail Modal**: Klik "Details" untuk melihat dashboard user lengkap
- **Statistics**: Summary cards dengan jumlah per status

## User Detail Modal Features

- **Basic Info**: User ID, Status, Session Code
- **Progress Info**: Untuk status IN_PROGRESS (jawaban, sisa waktu)  
- **Exam Results**: Untuk status COMPLETED (score, grade, hasil per kategori)
- **Real-time Data**: Fresh data dari API setiap kali modal dibuka

## Components Structure

```
src/
├── components/
│   ├── Layout.jsx          # Navigation & footer
│   └── UserDetailModal.jsx # Modal detail user
├── pages/
│   ├── Login.jsx           # Home page dengan User ID input
│   ├── ExamBoard.jsx       # Interface ujian
│   ├── Result.jsx          # Hasil ujian
│   └── AdminDashboard.jsx  # Dashboard admin
└── services/
    └── api.js              # API endpoints & axios config
```

## API Response Format

Semua API response menggunakan format:
```json
{
  "success": true/false,
  "message": "...",  
  "data": {...},
  "error": "..." // jika success=false
}
```

## Bootstrap Classes Used

- **Status Badges**: `bg-success`, `bg-primary`, `bg-warning`, `bg-danger`
- **Tables**: `table-striped`, `table-hover`, `table-responsive`
- **Cards**: `card`, `card-header`, `card-body`
- **Modals**: Bootstrap modal classes
- **Icons**: Bootstrap Icons (`bi-*`)

## Navigation

- **Navbar**: Home, Admin Dashboard, Options dropdown
- **Footer**: Simple footer dengan copyright
- **Active States**: Highlight halaman aktif di navbar