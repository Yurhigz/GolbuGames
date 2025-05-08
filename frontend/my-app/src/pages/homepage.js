import React from "react";
import { Link } from "react-router-dom";
import "./homepage.css";

const Home = () => {
    return (
        <div className="home-container">
            <div className="main-content-home">
                <div className="grid-preview">
                    {/* Exemple statique de grille */}
                    <div className="sudoku-preview">
                        {[...Array(9)].map((_, row) => (
                            <div key={row} className="preview-row">
                                {[...Array(9)].map((_, col) => (
                                    <div key={`${row}-${col}`} className="preview-cell">
                                        {/* Vous pouvez ajouter quelques chiffres fixes ici */}
                                        {Math.random() > 0.7 ? Math.floor(Math.random() * 9) + 1 : ""}
                                    </div>
                                ))}
                            </div>
                        ))}
                    </div>
                </div>
                <div className="welcome-section">
                    <h1>GolbuGames</h1>
                    <p className="subtitle">Jouer au Sudoku en Ligne</p>
                    <div className="game-modes">
                        <Link to="/multi" className="mode-button">Mode Multijoueur</Link>
                        <Link to="/solo" className="mode-button">Mode Solo</Link>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;