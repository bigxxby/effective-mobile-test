-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(255) NOT NULL UNIQUE,
    surname VARCHAR(255),
    name VARCHAR(255)
);

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL
);

-- Create task_logs table with ON DELETE CASCADE
CREATE TABLE IF NOT EXISTS task_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    task_id INT NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
