-- init-db.sql
DROP DATABASE IF EXISTS golbugamesdb;
CREATE DATABASE golbugamesdb;

\c golbugamesdb

-- Authentification générale peu importe le jeu
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    accountname VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL, --system de hash avant insertion
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gestion des sessions peu importe le jeu
DROP TABLE IF EXISTS sessions CASCADE;
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    game VARCHAR(100) NOT NULL,
    user_id INT NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Suppression en cascade en cas de suppresion d'un utilisateur, mécanisme à confirmer
);

DROP TABLE IF EXISTS sudoku_games CASCADE;
CREATE TABLE sudoku_games (
    id SERIAL PRIMARY KEY,
    board VARCHAR(200) NOT NULL,
    solution VARCHAR(200) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS games_scores CASCADE;
CREATE TABLE games_scores (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    opponent_id INT, -- optionnel selon le mode de jeu
    game_mode VARCHAR(20) NOT NULL CHECK (game_mode IN ('solo', '1v1')),
    results INT , -- 0 player 1 won, 1 even , 2 player 2 won
    completion_time INT,  -- temps en secondes lors des 1v1
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (opponent_id) REFERENCES users(id) ON DELETE SET NULL
);

DROP TABLE IF EXISTS user_stats CASCADE;
CREATE TABLE user_stats (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    total_games INT DEFAULT 0,
    total_wins INT DEFAULT 0,
    total_losses INT DEFAULT 0,
    total_draws INT DEFAULT 0,
    total_time INT DEFAULT 0, -- temps total de jeu en secondes
    average_time INT DEFAULT 0, -- temps moyen de jeu en secondes
    total_solo_games_finished INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS leaderboard CASCADE;
CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    elo_score INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Une ligne par relation d'ami
DROP TABLE IF EXISTS friendlist CASCADE;
CREATE TABLE friendlist (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    friend_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Une ligne par blocage d'utilisateur
DROP TABLE IF EXISTS blocked_users CASCADE;
CREATE TABLE blocked_users (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    blocked_user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS tournaments CASCADE;
CREATE TABLE tournaments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS tournament_participants CASCADE;
CREATE TABLE tournament_participants (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL,
    user_id INT NOT NULL,
    score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS cookies CASCADE;
CREATE TABLE cookies (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    cookie VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

)