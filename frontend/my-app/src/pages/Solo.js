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
  const [hints, setHints] = useState([]);
  const [selectedDifficulty, setSelectedDifficulty] = useState("Facile"); // ✅ UN SEUL état pour la difficulté
  const navigate = useNavigate();

  // Timer
  useEffect(() => {
    let interval;
    if (startTime) {
      interval = setInterval(() => {
        setElapsedTime(Math.floor((Date.now() - startTime) / 1000));
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [startTime]);

  // Quand une difficulté est choisie
  const handleSelectDifficulty = async (difficulty) => {
  setSelectedDifficulty(difficulty); // ✅ met à jour le state
  setIsModalOpen(false);
  const fetchedSolution = await fetchSolutionFromBackend(difficulty); 
  setSolution(fetchedSolution);
  setStartTime(Date.now());
  };

  // Récupération de la grille + indices simulés selon difficulté
  const fetchSolutionFromBackend = async (difficulty) => {
    try {
      // ⚠️ Pour l’instant statique → backend plus tard
      const solutionGrid = [
        [5,3,4,6,7,8,9,1,2],
        [6,7,2,1,9,5,3,4,8],
        [1,9,8,3,4,2,5,6,7],
        [8,5,9,7,6,1,4,2,3],
        [4,2,6,8,5,3,7,9,1],
        [7,1,3,9,2,4,8,5,6],
        [9,6,1,5,3,7,2,8,4],
        [2,8,7,4,1,9,6,3,5],
        [3,4,5,2,8,6,1,7,9],
      ];

      // ✅ Nombre d’indices selon la difficulté
      const generateHints = (difficulty) => {
        let hintsCount = 0;

        if (difficulty === "Facile") hintsCount = 20;
        if (difficulty === "Moyen") hintsCount = 15;
        if (difficulty === "Difficile") hintsCount = 10;

        const newHints = Array.from({ length: hintsCount }, (_, i) => `Indice ${i + 1}`);
        setHints(newHints);
        return hintsCount;
      };

      // ✅ Création du board avec indices aléatoires
      const hintsCount = generateHints(difficulty);
      const allIndexes = Array.from({ length: 81 }, (_, i) => i);

      const givenIndexes = new Set(
        allIndexes.sort(() => 0.5 - Math.random()).slice(0, hintsCount)
      );

      const board = solutionGrid.map((row, r) =>
        row.map((val, c) => (givenIndexes.has(r * 9 + c) ? val : null))
      );

      setGrid(board);
      setGivenCells(Array.from(givenIndexes));
      return solutionGrid;

    } catch (err) {
      console.error(err);
      return Array(9).fill(null).map(() => Array(9).fill(null));
    }
  };

  // Clic sur une case
  const handleCellClick = (row, col) => {
    if (givenCells.includes(row * 9 + col)) return;
    setSelectedCell({ row, col });
    if (selectedNumber !== null) {
      fillCell(row, col, selectedNumber);
    }
  };

  // Remplir une case
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

  // Grille complète ?
  const isGridComplete = () => grid.every(row => row.every(cell => cell !== null && cell !== ''));

  // Soumission de la grille
  const submitGrid = async () => {
    try {
      await fetch("http://localhost:8080/submit_solo_game", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          userId: 1,
          difficulty: selectedDifficulty, // ⚡ on garde la bonne difficulté
          time: elapsedTime,
          errors: errorCount,
          grid: grid.flat().join(""),
        }),
      });
    } catch (err) {
      console.error("Erreur lors de la soumission :", err);
    }
  };

  // Vérifie fin de partie
  useEffect(() => {
    if (solution && isGridComplete()) {
      alert(`✅ Grille terminée en ${elapsedTime} secondes avec ${errorCount} erreurs.`);
      submitGrid();
    }
  }, [grid, solution]);

  // Quitter
  const handleQuit = () => navigate('/');

  // Gestion clavier
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (isModalOpen) return;
      const { row, col } = selectedCell;
      if (givenCells.includes(row * 9 + col)) return;

      if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) {
        e.preventDefault();
      }
      if (e.key >= '1' && e.key <= '9') {
        fillCell(row, col, parseInt(e.key));
      } else if (e.key === 'Backspace' || e.key === 'Delete') {
        fillCell(row, col, null);
      } else if (e.key === 'ArrowUp' && row > 0) setSelectedCell({ row: row - 1, col });
      else if (e.key === 'ArrowDown' && row < 8) setSelectedCell({ row: row + 1, col });
      else if (e.key === 'ArrowLeft' && col > 0) setSelectedCell({ row, col: col - 1 });
      else if (e.key === 'ArrowRight' && col < 8) setSelectedCell({ row, col: col + 1 });
    };

    window.addEventListener('keydown', handleKeyDown, { passive: false });
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedCell, isModalOpen, grid, givenCells]);

  // Affichage grille
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
                    {cellValue}
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
        selectedDifficulty={selectedDifficulty} // ⚡ passe l’info au modal
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
