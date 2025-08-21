import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';

const Solo = () => {
    const [isModalOpen, setIsModalOpen] = useState(true);
    const [selectedNumber, setSelectedNumber] = useState(null);
    const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
    const [givenCells, setGivenCells] = useState([]);
    const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
    const [errorCount, setErrorCount] = useState(0);
    const [errorCells, setErrorCells] = useState([]);
    const [startTime, setStartTime] = useState(null);
    const [elapsedTime, setElapsedTime] = useState(0);
    const [selectedDifficulty, setSelectedDifficulty] = useState("easy");
    const navigate = useNavigate();

    // ========= Utils =========
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

    // Vérifie si une cellule est valide par rapport à toute la grille
    const isNumberValid = (row, col, number, currentGrid) => {
        if (!number) return true;

        // Vérifier ligne
        for (let c = 0; c < 9; c++) {
            if (c !== col && currentGrid[row][c] === number) return false;
        }

        // Vérifier colonne
        for (let r = 0; r < 9; r++) {
            if (r !== row && currentGrid[r][col] === number) return false;
        }

        // Vérifier carré 3x3
        const startRow = Math.floor(row / 3) * 3;
        const startCol = Math.floor(col / 3) * 3;
        for (let r = startRow; r < startRow + 3; r++) {
            for (let c = startCol; c < startCol + 3; c++) {
                if ((r !== row || c !== col) && currentGrid[r][c] === number) return false;
            }
        }

        return true;
    };

    // Recalcule toutes les erreurs de la grille
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

    const fillCell = (row, col, number) => {
        if (givenCells.includes(row * 9 + col)) return;
        const newGrid = grid.map(r => [...r]);
        newGrid[row][col] = number;

        // recalcul des erreurs AVANT la modif (ancien état)
        const prevErrors = recomputeErrors(grid);

        // application de la modif
        setGrid(newGrid);

        // recalcul des erreurs APRES la modif
        const newErrors = recomputeErrors(newGrid);
        setErrorCells(newErrors);

        // Vérifier si la nouvelle cellule est une erreur "nouvelle"
        if (number !== null && !isNumberValid(row, col, number, newGrid)) {
            const wasAlreadyInError = prevErrors.some(cell => cell.row === row && cell.col === col);
            if (!wasAlreadyInError) {
                setErrorCount((prev) => prev + 1);
            }
        }
    };

    // ========= Timer =========
    useEffect(() => {
        let interval;
        if (startTime) {
            interval = setInterval(() => {
                setElapsedTime(Math.floor((Date.now() - startTime) / 1000));
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [startTime]);

    // ========= Backend Helpers =========
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

            setIsModalOpen(false);
            setStartTime(Date.now());
        } catch (err) {
            if (err.response && err.response.status === 500) {
                await generateGridsIfNeeded(difficulty);
                return fetchGridFromBackend(difficulty);
            } else {
                console.error("Impossible de récupérer la grille :", err);
            }
        }
    };

    const handleSelectDifficulty = async (difficulty) => {
        setSelectedDifficulty(difficulty);
        await fetchGridFromBackend(difficulty);
    };

    // ========= Interactions =========
    const handleCellClick = (row, col) => {
        if (givenCells.includes(row * 9 + col)) return;
        setSelectedCell({ row, col });
        if (selectedNumber !== null) fillCell(row, col, selectedNumber);
    };

    const isGridComplete = () =>
        grid.every(row => row.every(cell => cell !== null && cell !== ''));

    const submitGrid = async () => {
        try {
            const gridString = grid.flat().map(v => v ?? 0).join('');
            await axios.post("http://localhost:3001/submit_solo_game", {
                userId: 1,
                difficulty: selectedDifficulty,
                time: elapsedTime,
                errors: errorCount,
                grid: gridString,
            });
        } catch (err) {
            console.error("Erreur lors de la soumission :", err);
        }
    };

    useEffect(() => {
        if (isGridComplete()) {
            alert(`✅ Grille terminée en ${elapsedTime} secondes avec ${errorCount} erreurs.`);
            submitGrid();
        }
    }, [grid]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleQuit = () => navigate('/');

    // ========= Clavier =========
    useEffect(() => {
        const handleKeyDown = (e) => {
            if (isModalOpen) return;
            const { row, col } = selectedCell;
            if (givenCells.includes(row * 9 + col)) return;

            if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) e.preventDefault();
            if (e.key >= '1' && e.key <= '9') fillCell(row, col, parseInt(e.key, 10));
            else if (e.key === 'Backspace' || e.key === 'Delete') fillCell(row, col, null);
            else if (e.key === 'ArrowUp' && row > 0) setSelectedCell({ row: row - 1, col });
            else if (e.key === 'ArrowDown' && row < 8) setSelectedCell({ row: row + 1, col });
            else if (e.key === 'ArrowLeft' && col > 0) setSelectedCell({ row, col: col - 1 });
            else if (e.key === 'ArrowRight' && col < 8) setSelectedCell({ row, col: col + 1 });
        };

        window.addEventListener('keydown', handleKeyDown, { passive: false });
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [selectedCell, isModalOpen, grid, givenCells]);

    // ========= UI =========
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

            {!isModalOpen && (
                <div className="game-content">
                    <SudokuGrid isBackground={false} />
                    <div className="actions-button-number">
                        <NumberSelector onNumberSelect={setSelectedNumber} selectedNumber={selectedNumber} />
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
