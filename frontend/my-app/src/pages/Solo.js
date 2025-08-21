import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';

const Solo = () => {
  const [isModalOpen, setIsModalOpen] = useState(true);
  const [selectedNumber, setSelectedNumber] = useState(null);
  const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
  const [givenCells, setGivenCells] = useState([]);
  const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
  const [solution, setSolution] = useState(null);
  const [errorCount, setErrorCount] = useState(0);
  const [errorCells, setErrorCells] = useState([]);
  const [startTime, setStartTime] = useState(null);
  const [elapsedTime, setElapsedTime] = useState(0);
  const [selectedDifficulty, setSelectedDifficulty] = useState("easy"); // "easy" | "intermediate" | "advanced" | "expert"
  const navigate = useNavigate();

  // ========= Utils =========

  // Convertit une string (81 chars) vers matrice 9x9 (null pour '.' ou '0')
  const parseBoard = (boardString) => {
    const cleaned = (boardString || "").replace(/[^0-9.]/g, "");
    if (cleaned.length !== 81) {
      console.warn("Board string length != 81:", cleaned.length, cleaned);
      return Array(9).fill(null).map(() => Array(9).fill(null));
    }
    return Array.from({ length: 9 }, (_, r) =>
      Array.from({ length: 9 }, (_, c) => {
        const ch = cleaned[r * 9 + c];
        return (ch === '.' || ch === '0') ? null : parseInt(ch, 10);
      })
    );
  };

  // Reconvertit une matrice en string (emptyChar = '.' ou '0')
  const stringifyBoard = (matrix, emptyChar = '0') => {
    return matrix
      .flat()
      .map(v => (v == null ? emptyChar : String(v)))
      .join('');
  };

  // Indices (0..80) des cases données (non nulles)
  const computeGivenIndexes = (matrix) => {
    const out = [];
    matrix.forEach((row, r) =>
      row.forEach((val, c) => {
        if (val != null) out.push(r * 9 + c);
      })
    );
    return out;
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

  // ========= Sélection difficulté =========
  const handleSelectDifficulty = async (difficulty) => {
    setSelectedDifficulty(difficulty);
    setIsModalOpen(false);

    await fetchGridFromBackend(difficulty);

    setStartTime(Date.now());
  };

  // ========= Backend: récupérer une grille =========
  const fetchGridFromBackend = async (difficulty) => {
    try {
      const res = await fetch(`http://localhost:3001/grid?difficulty=${difficulty}`);
      if (!res.ok) throw new Error(`Erreur API /grid: ${res.status}`);

      const data = await res.json();
      console.log("Grille récupérée :", data);

      const boardStr = data.board;        // obligatoire côté Go
      const solutionStr = data.solution;  // optionnel (selon ton handler)

      const board = parseBoard(boardStr);
      setGrid(board);
      setGivenCells(computeGivenIndexes(board));

      if (solutionStr) {
        setSolution(parseBoard(solutionStr));
      } else {
        setSolution(null); // si non fourni par l'API
      }
    } catch (err) {
      console.error("Impossible de récupérer la grille :", err);
    }
  };

  // ========= Interactions =========
  const handleCellClick = (row, col) => {
    if (givenCells.includes(row * 9 + col)) return; // lock cases données
    setSelectedCell({ row, col });
    if (selectedNumber !== null) {
      fillCell(row, col, selectedNumber);
    }
  };

  const fillCell = (row, col, number) => {
    if (givenCells.includes(row * 9 + col)) return;
    const newGrid = grid.map(r => [...r]);
    newGrid[row][col] = number;
    setGrid(newGrid);

    if (solution) {
      if (number === null) {
        setErrorCells(prev => prev.filter(cell => cell.row !== row || cell.col !== col));
      } else if (solution[row][col] !== number) {
        setErrorCount(prev => prev + 1);
        setErrorCells(prev => [...prev, { row, col }]);
      } else {
        setErrorCells(prev => prev.filter(cell => cell.row !== row || cell.col !== col));
      }
    }
  };

  const isGridComplete = () =>
    grid.every(row => row.every(cell => cell !== null && cell !== ''));

  // ========= Submit game =========
  const submitGrid = async () => {
    try {
      // si jamais la grille n'est pas totalement remplie, on envoie des '0'
      const gridString = stringifyBoard(grid, '0');

      await fetch("http://localhost:3001/submit_solo_game", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userId: 1,
          difficulty: selectedDifficulty,
          time: elapsedTime,
          errors: errorCount,
          grid: gridString,
        }),
      });
    } catch (err) {
      console.error("Erreur lors de la soumission :", err);
    }
  };

  // Fin de partie auto
  useEffect(() => {
    if (solution && isGridComplete()) {
      alert(`✅ Grille terminée en ${elapsedTime} secondes avec ${errorCount} erreurs.`);
      submitGrid();
    }
  }, [grid, solution]); // eslint-disable-line react-hooks/exhaustive-deps

  // Quitter
  const handleQuit = () => navigate('/');

  // ========= Clavier =========
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (isModalOpen) return;
      const { row, col } = selectedCell;
      if (givenCells.includes(row * 9 + col)) return;

      if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) {
        e.preventDefault();
      }
      if (e.key >= '1' && e.key <= '9') {
        fillCell(row, col, parseInt(e.key, 10));
      } else if (e.key === 'Backspace' || e.key === 'Delete') {
        fillCell(row, col, null);
      } else if (e.key === 'ArrowUp' && row > 0) setSelectedCell({ row: row - 1, col });
      else if (e.key === 'ArrowDown' && row < 8) setSelectedCell({ row: row + 1, col });
      else if (e.key === 'ArrowLeft' && col > 0) setSelectedCell({ row, col: col - 1 });
      else if (e.key === 'ArrowRight' && col < 8) setSelectedCell({ row, col: col + 1 });
    };

    window.addEventListener('keydown', handleKeyDown, { passive: false });
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedCell, isModalOpen, grid, givenCells]); // eslint-disable-line react-hooks/exhaustive-deps

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
      <div className="background-grid">
        <SudokuGrid isBackground={true} />
      </div>

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
