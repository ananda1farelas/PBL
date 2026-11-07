-- Data dosen pembimbing dan relasi bimbingan mahasiswa.
CREATE TABLE IF NOT EXISTS supervisors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    quota INTEGER NOT NULL CHECK (quota >= 0),
    supervised_count INTEGER NOT NULL DEFAULT 0 CHECK (supervised_count >= 0),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    CHECK (supervised_count <= quota)
);

CREATE TABLE IF NOT EXISTS supervisions (
    student_id INTEGER PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    supervisor_id INTEGER NOT NULL REFERENCES supervisors(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS supervisions_supervisor_id_idx ON supervisions(supervisor_id);

INSERT INTO permissions (name, description)
VALUES ('supervision:manage', 'Mengelola penugasan dosen pembimbing')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name)
VALUES ('admin', 'supervision:manage')
ON CONFLICT DO NOTHING;
