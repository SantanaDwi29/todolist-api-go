CREATE TABLE focus_sessions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    start_time DATETIME,
    end_time DATETIME,
    status ENUM('active', 'paused', 'completed') DEFAULT 'active',
    duration_minutes INT DEFAULT 45,
    paused_at DATETIME,
    elapsed_seconds INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
