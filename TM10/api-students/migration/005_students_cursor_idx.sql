-- Mendukung urutan keyset pagination pada GET /api/students.
-- id menjadi pemecah seri agar urutan tetap unik saat created_at sama.
CREATE INDEX IF NOT EXISTS students_created_at_id_desc_idx
    ON students (created_at DESC, id DESC);
