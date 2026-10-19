import React, { useCallback, useContext, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';
import { AuthContext } from "../contexts/AuthContext";
import { API_BASE_URL, WEBSOCKET_BASE_URL, authHeaders } from "../config/api";


const EndGameModal = ({ isOpen, onClose, points, time, errors, saveMessage }) => {
    if (!isOpen) return null;
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>🎉 Partie terminée !</h2>
                <p><strong>Points gagnés :</strong> {points}</p>
                <p><strong>Temps :</strong> {time} secondes</p>
                <p><strong>Erreurs :</strong> {errors}</p>
                {saveMessage && <p role="status">{saveMessage}</p>}
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

const calculatePoints = (difficulty, time, errors) => {
    const base = { easy: 100, intermediate: 200, hard: 300 }[difficulty] || 50;
    return Math.max(0, base - time - errors * 5);
};

const Solo = () => {

    const [isModalOpen, setIsModalOpen] = useState(true);
    const [showCountdown, setShowCountdown] = useState(false);
    const [selectedNumber, setSelectedNumber] = useState(null);
    const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
    const [givenCells, setGivenCells] = useState([]);
    const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
    const [highlightedNumber, setHighlightedNumber] = useState(null);
    const [errorCount, setErrorCount] = useState(0);
    const [errorCells, setErrorCells] = useState([]);
    const [startTime, setStartTime] = useState(null);
    const [elapsedTime, setElapsedTime] = useState(0);
    const [selectedDifficulty, setSelectedDifficulty] = useState("easy");
    const [showEndModal, setShowEndModal] = useState(false);
    const [saveMessage, setSaveMessage] = useState("");
    const [points, setPoints] = useState(0);
    const [isGridValid, setIsGridValid] = useState(false);
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();
    const socketRef = useRef(null);
    const completionSubmitted = useRef(false);

    // 🔥 Nouveau parseur pour tableau
    const parseBoardArray = (arr) => {
        if (!arr || !Array.isArray(arr) || arr.length !== 81) {
            return Array(9).fill(null).map(() => Array(9).fill(null));
        }
        const matrix = [];
        for (let i = 0; i < 9; i++) {
            matrix.push(
                arr.slice(i * 9, i * 9 + 9).map(num => (num === 0 ? null : num))
            );
        }
        return matrix;
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
        setHighlightedNumber(number);

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

    const handleSelectDifficulty = async (difficulty) => {
        setSelectedDifficulty(difficulty);
        setSaveMessage("");

        const socket = new WebSocket(`${WEBSOCKET_BASE_URL}/ws/solo?difficulty=${difficulty}`);
        socketRef.current = socket;

        socket.onopen = () => console.log("[WS] Connecté");

        socket.onmessage = (e) => {
            console.log("[WS] Reçu:", e.data);
            try {
                const msg = JSON.parse(e.data);

                if (msg.type === "init" && msg.grid) {
                    const board = parseBoardArray(msg.grid);
                    setGrid(board);
                    setGivenCells(computeGivenIndexes(board));
                    setShowCountdown(true);
                }
                else if (msg.position !== undefined && msg.value !== undefined && msg.valid !== undefined) {
                    const row = Math.floor(msg.position / 9);
                    const col = msg.position % 9;

                    if (msg.valid) {
                        setGivenCells(prev => [...prev, msg.position]);
                        setGrid(prevGrid => {
                            const newGrid = prevGrid.map(r => [...r]);
                            newGrid[row][col] = msg.value;
                            return newGrid;
                        });
                        setErrorCells(prev => prev.filter(c => !(c.row === row && c.col === col)));
                    } else {
                        setErrorCells(prev => [...prev, { row, col }]);
                        // setErrorCount(prev => prev + 1);
                    }
                }
            } catch (err) {
                console.error("[WS] Erreur parsing:", err);
            }
        };

        socket.onerror = (err) => console.error("[WS] Erreur:", err);
        socket.onclose = () => console.log("[WS] Fermé");
    };

    const startGameAfterCountdown = () => {
        setShowCountdown(false);
        setIsModalOpen(false);
        setStartTime(Date.now());
    };

    const submitGrid = useCallback(async (completionTime) => {
        if (!user) {
            setSaveMessage("Connectez-vous pour enregistrer votre résultat.");
            return;
        }
        try {
            await axios.post(`${API_BASE_URL}/submit_solo_game`, {
                user_id: user.id,
                difficulty: selectedDifficulty,
                completion_time: completionTime,
                game_mode: 'solo'
            }, {
                headers: authHeaders()
            });
            setSaveMessage("Résultat enregistré.");
        } catch {
            setSaveMessage("Impossible d’enregistrer votre résultat.");
        }
    }, [selectedDifficulty, user]);

    useEffect(() => {
        if (!isGridValid) {
            completionSubmitted.current = false;
            return;
        }
        if (completionSubmitted.current) return;

        completionSubmitted.current = true;
        const completionTime = startTime
            ? Math.max(1, Math.floor((Date.now() - startTime) / 1000))
            : 1;
        setElapsedTime(completionTime);
        setPoints(calculatePoints(selectedDifficulty, completionTime, errorCount));
        setShowEndModal(true);
        submitGrid(completionTime);
    }, [isGridValid, startTime, selectedDifficulty, errorCount, submitGrid]);

    const handleCellClick = (row, col) => {
        setSelectedCell({ row, col });
        setHighlightedNumber(grid[row][col]);
        if (!givenCells.includes(row * 9 + col) && selectedNumber !== null) {
            fillCell(row, col, selectedNumber);
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
            if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) e.preventDefault();
            if (e.key >= '1' && e.key <= '9') {
                if (!givenCells.includes(row * 9 + col)) fillCell(row, col, parseInt(e.key, 10));
            } else if (["Backspace", "Delete"].includes(e.key)) {
                if (!givenCells.includes(row * 9 + col)) fillCell(row, col, null);
            } else {
                let nextCell = null;
                if (e.key === 'ArrowUp' && row > 0) nextCell = { row: row - 1, col };
                else if (e.key === 'ArrowDown' && row < 8) nextCell = { row: row + 1, col };
                else if (e.key === 'ArrowLeft' && col > 0) nextCell = { row, col: col - 1 };
                else if (e.key === 'ArrowRight' && col < 8) nextCell = { row, col: col + 1 };
                if (nextCell) {
                    setSelectedCell(nextCell);
                    setHighlightedNumber(grid[nextCell.row][nextCell.col]);
                }
            }
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
                        const hasHighlightedNumber = highlightedNumber !== null && cellValue === highlightedNumber;
                        return (
                            <div
                                key={`${rowIndex}-${colIndex}`}
                                className={`grid-cell ${hasHighlightedNumber && !isBackground ? 'same-number-cell' : ''} ${isSelected && !isBackground ? 'selected-cell' : ''} ${isGiven ? 'given-cell' : ''}`}
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
                onClose={() => navigate("/")}
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
                saveMessage={saveMessage}
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
