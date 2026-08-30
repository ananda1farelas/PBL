Berikut adalah dokumentasi lengkap dan panduan langkah-demi-langkah yang siap kamu salin ke README repositori atau laporan tugas. Panduan ini ditulis secara mendetail agar rekan sekelas yang baru mengklona (*clone*) repositori bisa langsung menjalankan basis data dan aplikasi Go dari nol tanpa kebingungan.

# 🛠️ Panduan Konfigurasi Basis Data & Environment Setup

Dokumentasi ini berisi panduan untuk menyiapkan basis data PostgreSQL dari nol, skema tabel yang digunakan, serta daftar variabel environment (`.env`) yang diperlukan untuk menjalankan proyek backend ini.

## 1. Daftar Variabel Environment (`.env`)

Buat file baru bernama `.env` pada root direktori proyek (`api-students/.env`) dan isi dengan variabel berikut:

```env
# Server Configuration
PORT=8080

# PostgreSQL Database Configuration
APP_PORT=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=
DB_MAX_CONNS=

> **Catatan:** Sesuaikan `DB_PASSWORD` dengan password yang kamu tentukan saat membuat akun role PostgreSQL.

## 2. Skema Tabel (`students`)

Tabel `students` menyimpan data mahasiswa dengan skema sebagai berikut:

| Nama Kolom | Tipe Data | Aturan / Constraint | Deskripsi |
| --- | --- | --- | --- |
| `id` | `INTEGER` | `PRIMARY KEY`, `AUTO INCREMENT` | Identitas unik mahasiswa (otomatis) |
| `nim` | `VARCHAR(25)` | `NOT NULL`, `UNIQUE` | Nomor Induk Mahasiswa |
| `name` | `VARCHAR(50)` | `NOT NULL` | Nama lengkap mahasiswa |
| `grade` | `NUMERIC(5,2)` | `NOT NULL` | Nilai/IPK mahasiswa (contoh: 85.50) |
| `is_active` | `BOOLEAN` | `NOT NULL`, `DEFAULT true` | Status keaktifan mahasiswa |
| `created_at` | `TIMESTAMP` | `DEFAULT now()` | Waktu data dibuat |

## 3. Langkah-Langkah Menyiapkan Basis Data dari Nol

Ikuti urutan langkah berikut di terminal/PowerShell untuk menyiapkan PostgreSQL dan mengisi data awal:

 Langkah 1: Masuk ke PostgreSQL Superuser

Buka terminal dan masuk menggunakan user bawaan PostgreSQL:

powershell
psql -U postgres -d postgres

 Langkah 2: Buat User/Role dan Database

Eksekusi query SQL berikut untuk membuat role khusus proyek dan databasenya:

sql
 1. Buat role user baru
CREATE ROLE farelas WITH LOGIN SUPERUSER PASSWORD 'password_kamu_di_sini';

 2. Buat database baru
CREATE DATABASE praktikum_backend OWNER farelas;

 3. Keluar dari psql superuser
\q

 Langkah 3: Jalankan File Migrasi SQL

Masuk ke direktori proyek (`api-students`), lalu jalankan file migrasi untuk membuat tabel `students`:

powershell
psql -U farelas -d praktikum_backend -f migration/001_create_students.sql

*(Pastikan file `001_create_students.sql` berisi perintah DDL `CREATE TABLE` dan `CREATE INDEX` untuk tabel `students`)*.

 Langkah 4: Seeding Data Awal (5 Data Dummy)

Masuk ke prompt `psql` menggunakan user `farelas`:

powershell
psql -h 127.0.0.1 -U farelas -d praktikum_backend

Jalankan query `INSERT` berikut untuk memasukkan 5 data awal:

sql
INSERT INTO students (nim, name, grade, is_active) VALUES
('5025211001', 'Ahmad Dahlan', 85.50, true),
('5025211002', 'Budi Pratama', 90.00, true),
('5025211003', 'Citra Dewi', 78.25, true),
('5025211004', 'Deni Kurniawan', 88.75, false),
('5025211005', 'Eka Putri', 95.00, true);

 Langkah 5: Verifikasi Data

Pastikan data sudah masuk dengan benar menggunakan query:

sql
SELECT * FROM students;

Ketik `\q` untuk keluar dari psql.

## 4. Jalankan Aplikasi Go

Setelah file `.env` dibuat dan basis data siap, jalankan aplikasi Go dengan perintah:

powershell
go run main.go

Aplikasi sekarang siap menerima HTTP Request (seperti `POST`, `GET`) melalui port yang telah dikonfigurasi (`http://localhost:8080`).