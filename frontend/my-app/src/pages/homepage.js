import React, { useContext, useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { AuthContext } from "../contexts/AuthContext";
import "./homepage.css";

const Home = () => {
    const { user, isLoading } = useContext(AuthContext);
    const [grid, setGrid] = useState(
        Array(9).fill(null).map(() => Array(9).fill(null))
    );

    useEffect(() => {
        const interval = setInterval(() => {
            setGrid((prevGrid) => {
                const newGrid = prevGrid.map((row) => [...row]);
                const row = Math.floor(Math.random() * 9);
                const col = Math.floor(Math.random() * 9);

                newGrid[row][col] = newGrid[row][col]
                    ? null
                    : Math.floor(Math.random() * 9) + 1;

                return newGrid;
            });
        }, 300);

        return () => clearInterval(interval);
    }, []);

    return (
        <div className="home-container">
            <div className="main-content-home">
                <div className="grid-preview">
                    <div className="sudoku-preview">
                        {grid.map((row, rowIndex) => (
                            <div key={rowIndex} className="preview-row">
                                {row.map((cell, colIndex) => (
                                    <div key={`${rowIndex}-${colIndex}`} className="preview-cell">
                                        {cell || ""}
                                    </div>
                                ))}
                            </div>
                        ))}
                    </div>
                </div>
                <div className="welcome-section">
                    <h1>GolbuGames</h1>
                    <p className="subtitle">Jouer au Sudoku en Ligne</p>
                </div>
            </div>
            <div className="game-modes">
                {isLoading ? null : user ? (
                    <Link to="/multi" className="mode-button">Mode Multijoueur</Link>
                ) : (
                    <Link to="/login" className="mode-button">Connexion pour jouer en multijoueur</Link>
                )}
                <Link to="/solo" className="mode-button">Mode Solo</Link>
            </div>
        </div>
    );
};

export default Home;
