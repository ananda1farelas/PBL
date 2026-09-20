CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(25) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    grade DECIMAL(5,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nilai (
    id SERIAL PRIMARY KEY,
    id_student INT NOT NULL REFERENCES students(id),
    nama_mata_kuliah VARCHAR(50) NOT NULL,
    nilai DECIMAL(5,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS students_name_lower_idx 
    ON students (LOWER(name));