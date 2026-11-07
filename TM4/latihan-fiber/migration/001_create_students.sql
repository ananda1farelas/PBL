CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(25) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    grade DECIMAL(5,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS students_name_lower_idx 
    ON students (LOWER(name));