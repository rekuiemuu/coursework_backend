CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'doctor', 'technician')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    phone VARCHAR(20),
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_patients_last_name ON patients(last_name);
CREATE INDEX idx_patients_email ON patients(email);
CREATE INDEX idx_patients_phone ON patients(phone);

CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    device_id VARCHAR(255) NOT NULL UNIQUE,
    label VARCHAR(255),
    status VARCHAR(50) NOT NULL CHECK (status IN ('online', 'offline', 'error')),
    last_seen TIMESTAMP NOT NULL,
    brightness FLOAT NOT NULL DEFAULT 0.5 CHECK (brightness >= 0 AND brightness <= 1),
    zoom FLOAT NOT NULL DEFAULT 1.0 CHECK (zoom >= 1.0 AND zoom <= 10.0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_devices_device_id ON devices(device_id);
CREATE INDEX idx_devices_status ON devices(status);

CREATE TABLE IF NOT EXISTS examinations (
    id UUID PRIMARY KEY,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_examinations_patient_id ON examinations(patient_id);
CREATE INDEX idx_examinations_doctor_id ON examinations(doctor_id);
CREATE INDEX idx_examinations_status ON examinations(status);
CREATE INDEX idx_examinations_created_at ON examinations(created_at DESC);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL CHECK (file_size >= 0),
    mime_type VARCHAR(100) NOT NULL,
    width INTEGER NOT NULL CHECK (width >= 0),
    height INTEGER NOT NULL CHECK (height >= 0),
    captured_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_images_examination_id ON images(examination_id);
CREATE INDEX idx_images_filename ON images(filename);
CREATE INDEX idx_images_captured_at ON images(captured_at DESC);

CREATE TABLE IF NOT EXISTS analyses (
    id UUID PRIMARY KEY,
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    image_id UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    metrics JSONB,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_analyses_examination_id ON analyses(examination_id);
CREATE INDEX idx_analyses_image_id ON analyses(image_id);
CREATE INDEX idx_analyses_status ON analyses(status);
CREATE INDEX idx_analyses_metrics ON analyses USING GIN(metrics);

CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY,
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    diagnosis TEXT,
    recommendations TEXT,
    generated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_examination_id ON reports(examination_id);
CREATE INDEX idx_reports_generated_by ON reports(generated_by);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);

CREATE TABLE IF NOT EXISTS examination_images (
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    image_id UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    image_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (examination_id, image_id)
);

CREATE INDEX idx_examination_images_examination_id ON examination_images(examination_id);
CREATE INDEX idx_examination_images_image_id ON examination_images(image_id);
