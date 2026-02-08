-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    role INTEGER NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    action VARCHAR,
    device_name VARCHAR,
    ip_address VARCHAR,
    created_at TIMESTAMP
);

CREATE TABLE loan_slips (
    id SERIAL PRIMARY KEY,
    borrower_name VARCHAR,
    department VARCHAR,
    position VARCHAR,
    name VARCHAR,
    description TEXT,
    status VARCHAR,
    serial_number VARCHAR,
    images TEXT[],
    borrowed_date TIMESTAMP,
    returned_date TIMESTAMP,
    created_by INTEGER NOT NULL,
    updated_at TIMESTAMP,
    created_at TIMESTAMP
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    recipient_id INTEGER,
    sender_id INTEGER,
    title VARCHAR,
    type VARCHAR,
    content VARCHAR,
    is_read BOOLEAN,
    read_at TIMESTAMP,
    created_at TIMESTAMP
);

-- Foreign key constraints
ALTER TABLE loan_slips
    ADD CONSTRAINT fk_loan_slips_created_by
    FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE audit_logs
    ADD CONSTRAINT fk_audit_logs_user
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_recipient
    FOREIGN KEY (recipient_id) REFERENCES users(id);

ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_sender
    FOREIGN KEY (sender_id) REFERENCES users(id);

-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS loan_slips;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS users;
