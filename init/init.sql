-- init.sql

-- grant privileges to golang service user
GRANT ALL PRIVILEGES ON DATABASE messaging_db TO sample_user;

-- create table
CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    is_sent BOOLEAN NOT NULL DEFAULT FALSE
);

-- insert sample rows
INSERT INTO messages (id, content, phone_number) VALUES
    ('msg-001', 'Hello, world!', '+15551234567'),
    ('msg-002', 'Your order has shipped.', '+15557654321'),
    ('msg-003', 'Reminder: Your appointment is tomorrow.', '+15559876543'),
    ('msg-004', 'Welcome to our service!', '+15552345678'),
    ('msg-005', 'Your verification code is 4829.', '+15553456789'),
    ('msg-006', 'Happy birthday!', '+15554567890'),
    ('msg-007', 'Your subscription expires soon.', '+15555678901'),
    ('msg-008', 'New login from unknown device.', '+15556789012'),
    ('msg-009', 'Password reset request received.', '+15557890123'),
    ('msg-010', 'Thank you for your feedback.', '+15558901234'),
    ('msg-011', 'We have updated our terms of service.', '+15559012345'),
    ('msg-012', 'Your invoice is ready to view.', '+15550123456');
