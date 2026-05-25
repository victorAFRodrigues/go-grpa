CREATE TABLE IF NOT EXISTS config (
    id INT AUTO_INCREMENTPRIMARY KEY,
    grpa_name VARCHAR(100) NOT NULL,
    grpa_version VARCHAR(50) NOT NULL,
    worker_timeout INTEGER NOT NULL,
    system_name VARCHAR(100) NOT NULL,
    system_url VARCHAR(255) NOT NULL,
    system_username VARCHAR(100) NOT NULL,
    system_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);