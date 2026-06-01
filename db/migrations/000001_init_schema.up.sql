-- Migration UP digunakan untuk membuat schema database.
-- File ini dijalankan saat setup awal atau deploy ke server.

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'trainer', 'talent') NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_role (role)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS classes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS class_trainers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    class_id BIGINT UNSIGNED NOT NULL,
    trainer_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_class_trainer (class_id, trainer_id),
    CONSTRAINT fk_class_trainers_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_trainers_trainer FOREIGN KEY (trainer_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS class_talents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    class_id BIGINT UNSIGNED NOT NULL,
    talent_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_class_talent (class_id, talent_id),
    CONSTRAINT fk_class_talents_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_talents_talent FOREIGN KEY (talent_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS assignments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    class_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(150) NOT NULL,
    description TEXT NULL,
    deadline DATETIME NOT NULL,
    status ENUM('active', 'closed', 'archived') NOT NULL DEFAULT 'active',
    created_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_assignments_class_id (class_id),
    INDEX idx_assignments_deadline (deadline),
    CONSTRAINT fk_assignments_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignments_created_by FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS submissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    assignment_id BIGINT UNSIGNED NOT NULL,
    talent_id BIGINT UNSIGNED NOT NULL,
    github_url VARCHAR(255) NULL,
    deployment_url VARCHAR(255) NULL,
    notes TEXT NULL,
    status ENUM('not_submitted', 'submitted', 'reviewed', 'revision_required', 'late') NOT NULL DEFAULT 'submitted',
    submitted_at DATETIME NULL,
    reviewed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_submission_assignment_talent (assignment_id, talent_id),
    INDEX idx_submissions_status (status),
    INDEX idx_submissions_talent_id (talent_id),
    CONSTRAINT fk_submissions_assignment FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE,
    CONSTRAINT fk_submissions_talent FOREIGN KEY (talent_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS submission_reviews (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    submission_id BIGINT UNSIGNED NOT NULL,
    trainer_id BIGINT UNSIGNED NOT NULL,
    score INT NOT NULL,
    feedback TEXT NOT NULL,
    review_status ENUM('reviewed', 'revision_required') NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_reviews_submission_id (submission_id),
    CONSTRAINT fk_reviews_submission FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_trainer FOREIGN KEY (trainer_id) REFERENCES users(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS submission_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    submission_id BIGINT UNSIGNED NOT NULL,
    old_status VARCHAR(50) NOT NULL,
    new_status VARCHAR(50) NOT NULL,
    note TEXT NULL,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_histories_submission_id (submission_id),
    CONSTRAINT fk_histories_submission FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_histories_created_by FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB;

-- Password demo adalah SHA-256 dari: password123
-- Untuk production, gunakan bcrypt/argon2.

INSERT INTO users (id, name, email, password_hash, role) VALUES
(1, 'Admin Demo', 'admin@example.com', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 'admin'),
(2, 'Trainer Demo A', 'trainer@example.com', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 'trainer'),
(3, 'Trainer Demo B', 'trainer2@example.com', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 'trainer'),
(4, 'Talent Demo A', 'talent@example.com', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 'talent'),
(5, 'Talent Demo B', 'talent2@example.com', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', 'talent')
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    email = VALUES(email),
    password_hash = VALUES(password_hash),
    role = VALUES(role);

INSERT INTO classes (id, name, description, start_date, end_date) VALUES
(1, 'Golang Backend Batch 1', 'Class A untuk negative test object-level authorization', '2026-05-01', '2026-06-30'),
(2, 'Golang Backend Batch 2', 'Class B untuk negative test object-level authorization', '2026-06-01', '2026-07-31')
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    description = VALUES(description),
    start_date = VALUES(start_date),
    end_date = VALUES(end_date);

-- Bersihkan mapping salah dari seed lama.
DELETE FROM class_talents
WHERE class_id = 1 AND talent_id = 3;

-- Class A -> Trainer A
-- Class B -> Trainer B
INSERT INTO class_trainers (class_id, trainer_id) VALUES
(1, 2),
(2, 3)
ON DUPLICATE KEY UPDATE trainer_id = trainer_id;

-- Class A -> Talent A
-- Class B -> Talent B
INSERT INTO class_talents (class_id, talent_id) VALUES
(1, 4),
(2, 5)
ON DUPLICATE KEY UPDATE talent_id = talent_id;

-- Assignment A -> Class A
-- Assignment B -> Class B
INSERT INTO assignments (id, class_id, title, description, deadline, status, created_by) VALUES
(1, 1, 'Assignment A - REST API CRUD Product', 'Assignment milik Class A', '2026-06-30 23:59:00', 'active', 2),
(2, 2, 'Assignment B - REST API CRUD Order', 'Assignment milik Class B', '2026-07-31 23:59:00', 'active', 3)
ON DUPLICATE KEY UPDATE
    class_id = VALUES(class_id),
    title = VALUES(title),
    description = VALUES(description),
    deadline = VALUES(deadline),
    status = VALUES(status),
    created_by = VALUES(created_by);

-- Submission A -> milik Talent A untuk Assignment A
-- Submission B -> milik Talent B untuk Assignment B
INSERT INTO submissions (
    id,
    assignment_id,
    talent_id,
    github_url,
    deployment_url,
    notes,
    status,
    submitted_at
) VALUES
(1, 1, 4, 'https://github.com/demo/talent-a-assignment-a', 'https://talent-a-demo.example.com', 'Submission A milik Talent A', 'submitted', '2026-06-10 10:00:00'),
(2, 2, 5, 'https://github.com/demo/talent-b-assignment-b', 'https://talent-b-demo.example.com', 'Submission B milik Talent B', 'submitted', '2026-06-11 10:00:00')
ON DUPLICATE KEY UPDATE
    assignment_id = VALUES(assignment_id),
    talent_id = VALUES(talent_id),
    github_url = VALUES(github_url),
    deployment_url = VALUES(deployment_url),
    notes = VALUES(notes),
    status = VALUES(status),
    submitted_at = VALUES(submitted_at);