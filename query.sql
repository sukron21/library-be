 -- Membuat database baru
CREATE DATABASE library_db;

-- Tabel books
CREATE TABLE books (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0),
    category VARCHAR(100)
);

-- Tabel users (ganti dari borrowers)
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE
);

-- Tabel lending_records
CREATE TABLE lending_records (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL,
    user_id UUID NOT NULL,        -- ganti borrower_id jadi user_id
    borrow_date DATE NOT NULL,
    return_date DATE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
