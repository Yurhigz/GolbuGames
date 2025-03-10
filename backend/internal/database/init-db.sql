-- init-db.sql
CREATE DATABASE sudokuDB;

\c sudokuDB

CREATE TABLE sudoku (
    id SERIAL PRIMARY KEY,
    board VARCHAR(81) NOT NULL,
    solution VARCHAR(81) NOT NULL,
    difficulty VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL, --system de hash avant insertion
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Suppression en cascade en cas de suppresion d'un utilisateur, mécanisme à confirmer
);

CREATE TABLE scores (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);

CREATE DATABASE membersGG;