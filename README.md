# Sistem Manajemen SDM (HRMS) - Backend

Backend Sistem Manajemen Sumber Daya Manusia (HRMS) yang komprehensif, dibangun menggunakan Go, framework Gin, dan PostgreSQL. Proyek ini dirancang untuk menangani berbagai aspek manajemen karyawan secara efisien dan aman.

## 🌟 Fitur Utama

- **Otentikasi & Otorisasi**: Sistem login aman berbasis JWT dengan kontrol akses berbasis peran (RBAC).
- **Manajemen Karyawan**: Operasi CRUD lengkap untuk data karyawan, termasuk profil dan riwayat.
- **Pelacakan Kehadiran**: Sistem clock-in/clock-out real-time dengan perhitungan jam kerja otomatis.
- **Manajemen Cuti**: Pengajuan cuti, alur persetujuan (approval workflow), dan pelacakan sisa cuti.
- **Pemrosesan Penggajian**: Pembuatan slip gaji bulanan dan pemrosesan pembayaran yang akurat.
- **Manajemen Departemen**: Struktur hierarki departemen dan penugasan karyawan.

## 🛠️ Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go 1.21+
- **Web Framework**: Gin (High-performance HTTP web framework)
- **ORM**: GORM (The fantastic ORM library for Golang)
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Otentikasi**: JWT (JSON Web Tokens)

## 📋 Prasyarat

Sebelum memulai, pastikan Anda telah menginstal:

- Go 1.21 atau lebih baru
- PostgreSQL 15 atau lebih baru
- Redis 7 atau lebih baru (opsional, untuk caching)

## 🚀 Instalasi & Menjalankan Aplikasi

1. **Clone repositori:**
   ```bash
   git clone <repository-url>
   cd hr-backend
   ```

2. **Instal dependensi:**
   ```bash
   go mod download
   ```

3. **Konfigurasi Environment:**
   Salin file contoh `.env` dan sesuaikan dengan konfigurasi lokal Anda.
   ```bash
   cp .env.example .env
   ```

4. **Jalankan Aplikasi:**
   Anda dapat menjalankan aplikasi langsung menggunakan Go:
   ```bash
   go run cmd/api/main.go
   ```
   Atau menggunakan Makefile:
   ```bash
   make run
   ```

## 🗄️ Pengaturan Database

**Menggunakan Docker (Direkomendasikan):**
```bash
make db-create
make db-start
```

**Secara Manual:**
Masuk ke PostgreSQL dan jalankan perintah berikut:
```sql
CREATE DATABASE hrms_db;
CREATE USER hrms_user WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE hrms_db TO hrms_user;
```

## 🔌 Endpoint API Utama

Berikut adalah beberapa endpoint kunci yang tersedia:

### Otentikasi
- `POST /api/v1/auth/login` - Masuk ke sistem
- `POST /api/v1/auth/logout` - Keluar dari sistem
- `GET /api/v1/auth/profile` - Lihat profil pengguna

### Karyawan (Employees)
- `GET /api/v1/employees` - Daftar semua karyawan
- `POST /api/v1/employees` - Tambah karyawan baru
- `GET /api/v1/employees/:id` - Detail karyawan

### Kehadiran (Attendance)
- `POST /api/v1/attendance/clock-in` - Catat jam masuk
- `POST /api/v1/attendance/clock-out` - Catat jam keluar
- `GET /api/v1/attendance/report` - Laporan kehadiran

### Cuti (Leaves)
- `POST /api/v1/leaves` - Ajukan permohonan cuti
- `PUT /api/v1/leaves/:id/approve` - Setujui/Tolak cuti

### Penggajian (Payroll)
- `POST /api/v1/payroll/generate` - Generate gaji bulanan
- `POST /api/v1/payroll/:id/process-payment` - Proses pembayaran gaji

## 👤 User Admin Default

Setelah menjalankan migrasi database, Anda dapat menggunakan akun admin berikut untuk pengujian:

- **Email**: `admin@hrms.com`
- **Password**: `admin123`
- **Role**: `admin`

## 🐳 Dukungan Docker

Untuk menjalankan seluruh aplikasi menggunakan Docker:

**Build Image:**
```bash
make docker-build
```

**Jalankan Container:**
```bash
make docker-run
```

## 🧪 Pengujian (Testing)

Jalankan unit test untuk memastikan semua fungsi berjalan dengan baik:
```bash
make test
```

## 📂 Struktur Proyek

Proyek ini mengikuti prinsip **Clean Architecture** untuk memastikan kode yang modular dan mudah dipelihara:

```
hr-backend/
├── cmd/
│   └── api/          # Entry point aplikasi
├── internal/
│   ├── config/       # Konfigurasi aplikasi
│   ├── database/     # Setup koneksi database
│   ├── handlers/     # HTTP Handlers (Controllers)
│   ├── middleware/   # Middleware (Auth, Logger, dll)
│   ├── models/       # Definisi struct & model database
│   ├── repositories/ # Data Access Layer
│   ├── services/     # Business Logic Layer
│   └── utils/        # Fungsi utilitas
├── .env.example      # Template variabel environment
├── Dockerfile        # Konfigurasi Docker
├── Makefile          # Shortcut perintah
└── README.md         # Dokumentasi proyek
```

## 📄 Lisensi

MIT License
