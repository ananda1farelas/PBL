# Student Management System - RESTful API

RESTful API untuk pengelolaan data mahasiswa (*Student Management System*) yang dibangun menggunakan bahasa pemrosesan **Go (Golang)** dan *framework* **Fiber**. API ini dirancang mengikuti standar arsitektur REST, dilengkapi validasi DTO ketat, *custom middleware*, penanganan kesalahan tersentralisasi, serta standarisasi format respons JSON.

---

## 🚀 Fitur Utama

- **Layanan CRUD Lengkap**: Operasi `GET` (list & detail), `POST`, `PUT` (full replacement), `PATCH` (partial update), dan `DELETE`.
- **Query String Lanjutan**:
  - Paginasi data (`page` & `limit`) dengan pembatasan *max limit* (proteksi OOM/DoS).
  - Pencarian nama *case-insensitive* (`search`).
  - Pengurutan data fleksibel berbasis *whitelist* field (`sort` & `order`).
  - Penyaringan data berbasis kriteria (`is_active`).
- **Keamanan & Middleware**:
  - Middleware `requireJSON` untuk validasi header `Content-Type: application/json` (`415 Unsupported Media Type`).
  - Validasi *input body* menggunakan *struct tags* DTO (`422 Unprocessable Entity`).
  - Deteksi konflik duplikasi NIM (`409 Conflict`).
- **Konsistensi Respons**: Menggunakan struktur *Standard JSON Envelope* untuk semua status respons (sukses maupun error).

---

## 🛠️ Prasyarat & Instalasi

Pastikan **Go (v1.20+)** telah terinstal di perangkat kamu.

1. **Clone repositori ini**:
   ```bash
   git clone [https://github.com/username/student-management-api.git](https://github.com/username/student-management-api.git)
   cd student-management-api
   