-- init-db.sql
CREATE DATABASE golbugamesdb;

\c golbugamesdb

-- Authentification générale peu importe le jeu
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    accountname VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL, --system de hash avant insertion
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gestion des sessions peu importe le jeu 
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    game VARCHAR(100) NOT NULL,
    user_id INT NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Suppression en cascade en cas de suppresion d'un utilisateur, mécanisme à confirmer
);


CREATE TABLE sudoku_games (
    id SERIAL PRIMARY KEY,
    board VARCHAR(160) NOT NULL,
    solution VARCHAR(160) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE games_scores (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    game_type VARCHAR(100) NOT NULL,
    score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);