import React, { useContext, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';
import { AuthContext } from "../contexts/AuthContext";

const EndGameModal = ({ isOpen, onClose, points, time, errors }) => {
    if (!isOpen) return null;
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>🎉 Partie terminée !</h2>
                <p><strong>Points gagnés :</strong> {points}</p>
                <p><strong>Temps :</strong> {time} secondes</p>
                <p><strong>Erreurs :</strong> {errors}</p>
                <button onClick={onClose} className="modal-button">OK</button>
            </div>
        </div>
    );
};

const CountdownOverlay = ({ onComplete }) => {
    const [count, setCount] = useState(3);

    useEffect(() => {
        if (count === 0) {
            const timer = setTimeout(() => onComplete(), 500);
            return () => clearTimeout(timer);
        }
        const interval = setInterval(() => setCount((prev) => prev - 1), 1000);
        return () => clearInterval(interval);
    }, [count, onComplete]);

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h1 style={{ fontSize: '3rem' }}>{count === 0 ? 'Partez !' : count}</h1>
            </div>
        </div>
    );
};

const Solo = () => {
    const [isModalOpen, setIsModalOpen] = useState(true);
    const [showCountdown, setShowCountdown] = useState(false);
    const [selectedNumber, setSelectedNumber] = useState(null);
    const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
    const [givenCells, setGivenCells] = useState([]);
    const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
    const [errorCount, setErrorCount] = useState(0);
    const [errorCells, setErrorCells] = useState([]);
    const [startTime, setStartTime] = useState(null);
    const [elapsedTime, setElapsedTime] = useState(0);
    const [selectedDifficulty, setSelectedDifficulty] = useState("easy");
    const [showEndModal, setShowEndModal] = useState(false);
    const [points, setPoints] = useState(0);
    const [isGridValid, setIsGridValid] = useState(false);
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();
    const socketRef = useRef(null);

    const parseBoardString = (boardStr) => {
        if (!boardStr) return Array(9).fill(null).map(() => Array(9).fill(null));
        return boardStr
            .split(";")
            .filter(line => line.trim() !== "")
            .map(row =>
                row.split(",").map(val => {
                    const num = parseInt(val, 10);
                    return isNaN(num) || num === 0 ? null : num;
                })
            );
    };

    const computeGivenIndexes = (matrix) => {
        const out = [];
        matrix.forEach((row, r) =>
            row.forEach((val, c) => {
                if (val != null) out.push(r * 9 + c);
            })
        );
        return out;
    };

    const isNumberValid = (row, col, number, currentGrid) => {
        if (!number) return true;
        for (let c = 0; c < 9; c++) if (c !== col && currentGrid[row][c] === number) return false;
        for (let r = 0; r < 9; r++) if (r !== row && currentGrid[r][col] === number) return false;
        const startRow = Math.floor(row / 3) * 3;
        const startCol = Math.floor(col / 3) * 3;
        for (let r = startRow; r < startRow + 3; r++) {
            for (let c = startCol; c < startCol + 3; c++) {
                if ((r !== row || c !== col) && currentGrid[r][c] === number) return false;
            }
        }
        return true;
    };

    const recomputeErrors = (currentGrid) => {
        const errors = [];
        currentGrid.forEach((row, r) => {
            row.forEach((val, c) => {
                if (val !== null && !isNumberValid(r, c, val, currentGrid)) {
                    errors.push({ row: r, col: c });
                }
            });
        });
        return errors;
    };

    const isGridComplete = (currentGrid) =>
        currentGrid.every(row => row.every(cell => cell !== null && cell !== ''));

    const isGridFullyValid = (currentGrid) =>
        isGridComplete(currentGrid) &&
        currentGrid.every((row, r) =>
            row.every((val, c) => isNumberValid(r, c, val, currentGrid))
        );

    const fillCell = (row, col, number) => {
        if (givenCells.includes(row * 9 + col)) return;
        const newGrid = grid.map(r => [...r]);
        newGrid[row][col] = number;

        const prevErrors = recomputeErrors(grid);
        setGrid(newGrid);

        const newErrors = recomputeErrors(newGrid);
        setErrorCells(newErrors);

        if (number !== null && !isNumberValid(row, col, number, newGrid)) {
            const wasAlreadyInError = prevErrors.some(cell => cell.row === row && cell.col === col);
            if (!wasAlreadyInError) {
                setErrorCount((prev) => prev + 1);
            }
        }

        // Envoi WebSocket
        const position = row * 9 + col;
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify({ position, value: number }));
        }

        setIsGridValid(isGridFullyValid(newGrid));
    };

    useEffect(() => {
        let interval;
        if (startTime && !showEndModal) {
            interval = setInterval(() => {
                setElapsedTime(Math.floor((Date.now() - startTime) / 1000));
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [startTime, showEndModal]);

    const generateGridsIfNeeded = async (difficulty) => {
        try {
            const promises = Array.from({ length: 10 }).map(() =>
                axios.post("http://127.0.0.1:3001/add_grid", { Difficulty: difficulty })
            );
            await Promise.all(promises);
        } catch (err) {
            console.error("Impossible de générer les grilles :", err);
        }
    };

    const fetchGridFromBackend = async (difficulty) => {
        try {
            const res = await axios.get(`http://127.0.0.1:3001/grid?difficulty=${difficulty}`);
            const data = res.data;
            const board = parseBoardString(data.board);
            setGrid(board);
            setGivenCells(computeGivenIndexes(board));
            setShowCountdown(true);
        } catch (err) {
            if (err.response?.status === 500) {
                await generateGridsIfNeeded(difficulty);
                return fetchGridFromBackend(difficulty);
            } else {
                console.error("Impossible de récupérer la grille :", err);
            }
        }
    };

    const handleSelectDifficulty = async (difficulty) => {
        setSelectedDifficulty(difficulty);

        // Connexion WebSocket
        const socket = new WebSocket("ws://localhost:3001/ws/solo");
        socketRef.current = socket;

        socket.onopen = () => console.log("[WS] Connecté");
        socket.onmessage = (e) => console.log("[WS] Reçu:", e.data);
        socket.onerror = (err) => console.error("[WS] Erreur:", err);
        socket.onclose = () => console.log("[WS] Fermé");

        await fetchGridFromBackend(difficulty);
    };

    const startGameAfterCountdown = () => {
        setShowCountdown(false);
        setIsModalOpen(false);
        setStartTime(Date.now());
    };

    const calculatePoints = (difficulty, time, errors) => {
        const base = { easy: 100, intermediate: 200, advanced: 300, expert: 500 }[difficulty] || 50;
        let score = base - time - errors * 5;
        return score > 0 ? score : 0;
    };

    const submitGrid = async () => {
        try {
            await axios.post("http://localhost:3001/submit_solo_game", {
                user_id: user.id,
                difficulty: selectedDifficulty,
                completion_time: elapsedTime,
                game_mode: 'sudoku'
            });
        } catch (err) {
            console.error("Erreur soumission:", err);
        }
    };

    useEffect(() => {
        if (isGridValid) {
            const pts = calculatePoints(selectedDifficulty, elapsedTime, errorCount);
            setPoints(pts);
            setShowEndModal(true);
            submitGrid();
        }
    }, [isGridValid]);

    const handleCellClick = (row, col) => {
        if (!givenCells.includes(row * 9 + col)) {
            setSelectedCell({ row, col });
            if (selectedNumber !== null) {
                fillCell(row, col, selectedNumber);
            }
        }
    };

    const handleNumberSelect = (number) => {
        setSelectedNumber(prev => (prev === number ? null : number));
    };

    const handleQuit = () => {
        if (socketRef.current) socketRef.current.close();
        navigate('/');
    };

    useEffect(() => {
        const handleKeyDown = (e) => {
            if (isModalOpen) return;
            const { row, col } = selectedCell;
            if (givenCells.includes(row * 9 + col)) return;
            if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) e.preventDefault();
            if (e.key >= '1' && e.key <= '9') fillCell(row, col, parseInt(e.key, 10));
            else if (["Backspace", "Delete"].includes(e.key)) fillCell(row, col, null);
            else if (e.key === 'ArrowUp' && row > 0) setSelectedCell({ row: row - 1, col });
            else if (e.key === 'ArrowDown' && row < 8) setSelectedCell({ row: row + 1, col });
            else if (e.key === 'ArrowLeft' && col > 0) setSelectedCell({ row, col: col - 1 });
            else if (e.key === 'ArrowRight' && col < 8) setSelectedCell({ row, col: col + 1 });
        };
        window.addEventListener('keydown', handleKeyDown, { passive: false });
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [selectedCell, isModalOpen, grid, givenCells]);

    const SudokuGrid = ({ isBackground }) => (
        <div className={`sudoku-grid ${isBackground ? 'background' : 'active'}`}>
            {[...Array(9)].map((_, rowIndex) => (
                <div key={rowIndex} className="grid-row">
                    {[...Array(9)].map((_, colIndex) => {
                        const isSelected = selectedCell.row === rowIndex && selectedCell.col === colIndex;
                        const hasError = errorCells.some(cell => cell.row === rowIndex && cell.col === colIndex);
                        const cellValue = grid[rowIndex][colIndex];
                        const isGiven = givenCells.includes(rowIndex * 9 + colIndex);
                        return (
                            <div
                                key={`${rowIndex}-${colIndex}`}
                                className={`grid-cell ${isSelected && !isBackground ? 'selected-cell' : ''} ${isGiven ? 'given-cell' : ''}`}
                                onClick={!isBackground ? () => handleCellClick(rowIndex, colIndex) : undefined}
                            >
                                {!isBackground && (
                                    <span style={{ color: hasError ? 'red' : isGiven ? 'blue' : 'black' }}>
                                        {cellValue ?? ""}
                                    </span>
                                )}
                            </div>
                        );
                    })}
                </div>
            ))}
        </div>
    );

    return (
        <div className="solo-container">
            <div className="background-grid"><SudokuGrid isBackground={true} /></div>

            <DifficultyModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSelectDifficulty={handleSelectDifficulty}
                selectedDifficulty={selectedDifficulty}
                setSelectedDifficulty={setSelectedDifficulty}
            />

            {showCountdown && <CountdownOverlay onComplete={startGameAfterCountdown} />}

            <EndGameModal
                isOpen={showEndModal}
                onClose={() => setShowEndModal(false)}
                points={points}
                time={elapsedTime}
                errors={errorCount}
            />

            {!isModalOpen && !showCountdown && (
                <div className="game-content">
                    <SudokuGrid isBackground={false} />
                    <div className="actions-button-number">
                        <NumberSelector
                            onNumberSelect={handleNumberSelect}
                            selectedNumber={selectedNumber}
                        />
                        <div className="game-actions">
                            <button className="quit-button" onClick={handleQuit}>Quitter</button>
                            <div className="error-counter">Nombre d'erreurs : {errorCount}</div>
                            <div className="timer">Temps : {elapsedTime}s</div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Solo;
