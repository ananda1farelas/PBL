# Modul 9 - D2: Pembimbing Tugas Akhir untuk Entity Students

Implementasi ini melanjutkan project Modul 7 `api-students`.

## File baru
- `app/model/supervision.go`
- `app/repository/tx.go`
- `app/repository/supervision_repository.go`
- `app/service/supervision_service.go`
- `migration/006_create_supervisions.sql`

## File yang diubah
- `app/repository/students_repository.go` (menggunakan DBTX dan WithTx)
- `middleware/middleware.go` (permission middleware)
- `helper/errors.go` (401 dan 403)
- `route/route.go`
- `config/app.go`
- `main.go`

## Menjalankan migration
Jalankan seluruh migration secara berurutan, termasuk:

```sql
\i migration/006_create_supervisions.sql
```

Pastikan migration 001 sampai 005 sudah dijalankan terlebih dahulu karena migration 006 menggunakan tabel `students`, `permissions`, dan `role_permissions`.

## Contoh data dosen pembimbing
Setelah migration selesai, tambahkan data awal melalui PostgreSQL:

```sql
INSERT INTO supervisors (name, quota) VALUES
('Dr. Dosen Satu', 3),
('Dr. Dosen Dua', 5);
```

## Endpoint
Semua endpoint berikut memerlukan access token JWT dengan role `admin` (permission `supervision:manage`).

### Assign pembimbing
`POST /api/supervisors/{id}/assign`

```json
{"student_id": 1}
```

### Pindah pembimbing
`POST /api/supervisions/switch`

```json
{"student_id": 1, "new_supervisor_id": 2}
```

### Update data/kuota pembimbing (optimistic locking)
`PATCH /api/supervisors/{id}`

```json
{"quota": 4, "expected_version": 1}
```

`name` bersifat opsional. `expected_version` wajib dikirim sesuai nilai version terakhir yang dibaca client. Jika version sudah berubah, API mengembalikan conflict agar client memuat ulang data.

## Catatan concurrency
- `supervisions.student_id` adalah primary key, sehingga satu mahasiswa hanya dapat memiliki satu pembimbing.
- Assign mengunci baris mahasiswa dengan `SELECT ... FOR UPDATE` dan menaikkan `supervised_count` melalui conditional update `supervised_count < quota`.
- Switch berjalan dalam satu transaksi. Baris mahasiswa dikunci lebih dahulu; dua baris supervisor dikunci dalam urutan ID menaik. Jika kuota pembimbing baru penuh, transaksi rollback dan relasi pembimbing lama tetap utuh.
- Update quota menggunakan `WHERE id = ? AND version = ?` dan menaikkan version setiap kali update berhasil.
- Transaksi dibuat melalui `TxManager.WithTx`; repository dapat dibungkus dengan `WithTx` untuk memakai koneksi transaksi yang sama.

Permission `supervision:manage` ditambahkan ke tabel permission dan diberikan kepada role `admin` melalui migration. Pemeriksaan route pada implementasi ini memetakan permission tersebut ke role `admin`.
