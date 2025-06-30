import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';

const Solo = () => {
  // États pour gérer les modales, la difficulté, le numéro sélectionné, la grille de jeu et la cellule sélectionnée
  const [isModalOpen, setIsModalOpen] = useState(true);
  const [difficulty, setDifficulty] = useState(null);
  const [selectedNumber, setSelectedNumber] = useState(null);
  const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
  const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
  const [solution, setSolution] = useState(null);
  const [errorCount, setErrorCount] = useState(0);
  const [errorCells, setErrorCells] = useState([]);
  const [startTime, setStartTime] = useState(null);
  const [elapsedTime, setElapsedTime] = useState(0);
  const navigate = useNavigate();

  // Timer qui met à jour le temps écoulé
  useEffect(() => {
    let interval;
    if (startTime) {
      interval = setInterval(() => {
        setElapsedTime(Math.floor((Date.now() - startTime) / 1000));
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [startTime]);

  // Lorsque l'utilisateur choisit une difficulté, on charge la solution et on démarre le timer
  const handleSelectDifficulty = async (selectedDifficulty) => {
    setDifficulty(selectedDifficulty);
    setIsModalOpen(false);
    const fetchedSolution = await fetchSolutionFromBackend(selectedDifficulty);
    setSolution(fetchedSolution);
    setStartTime(Date.now());
  };

  // Simule la récupération d'une solution depuis le backend
  const fetchSolutionFromBackend = async (difficulty) => {
    return Array(9).fill([1,2,3,4,5,6,7,8,9]);
  };

  // Quand l'utilisateur clique sur une case
  const handleCellClick = (row, col) => {
    setSelectedCell({ row, col });
    if (selectedNumber !== null) {
      fillCell(row, col, selectedNumber);
    }
  };

  // Remplir une case avec un numéro, vérifier les erreurs
  const fillCell = (row, col, number) => {
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

  // Vérifie si la grille est complète
  const isGridComplete = () => {
    return grid.every(row => row.every(cell => cell !== null && cell !== ''));
  };

  // Si la grille est complète, afficher le message de succès avec stats
  useEffect(() => {
    if (solution && isGridComplete()) {
      alert(`✅ Grille terminée en ${elapsedTime} secondes avec ${errorCount} erreurs.`);
    }
  }, [grid, solution]);

  // Quitter la partie
  const handleQuit = () => {
    navigate('/');
  };

  // Gestion des touches du clavier (numéros et déplacements)
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (isModalOpen) return;
      if (["ArrowUp","ArrowDown","ArrowLeft","ArrowRight"].includes(e.key)) {
        e.preventDefault(); // Empêche l'ascenseur de bouger
      }
      const { row, col } = selectedCell;
      if (e.key >= '1' && e.key <= '9') {
        fillCell(row, col, parseInt(e.key));
      } else if (e.key === 'Backspace' || e.key === 'Delete') {
        fillCell(row, col, null);
      } else if (e.key === 'ArrowUp' && row > 0) {
        setSelectedCell({ row: row - 1, col });
      } else if (e.key === 'ArrowDown' && row < 8) {
        setSelectedCell({ row: row + 1, col });
      } else if (e.key === 'ArrowLeft' && col > 0) {
        setSelectedCell({ row, col: col - 1 });
      } else if (e.key === 'ArrowRight' && col < 8) {
        setSelectedCell({ row, col: col + 1 });
      }
    };

    window.addEventListener('keydown', handleKeyDown, { passive: false });
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedCell, isModalOpen, grid]);

  // Composant pour afficher la grille de Sudoku
  const SudokuGrid = ({ isBackground }) => (
    <div className={`sudoku-grid ${isBackground ? 'background' : 'active'}`}>
      {[...Array(9)].map((_, rowIndex) => (
        <div key={rowIndex} className="grid-row">
          {[...Array(9)].map((_, colIndex) => {
            const isSelected = selectedCell.row === rowIndex && selectedCell.col === colIndex;
            const hasError = errorCells.some(cell => cell.row === rowIndex && cell.col === colIndex);
            const cellValue = grid[rowIndex][colIndex];
            return (
              <div
                key={`${rowIndex}-${colIndex}`}
                className={`grid-cell ${isSelected && !isBackground ? 'selected-cell' : ''}`}
                onClick={!isBackground ? () => handleCellClick(rowIndex, colIndex) : undefined}
              >
                {!isBackground && (
                  <span style={{ color: hasError ? 'red' : 'black' }}>{cellValue}</span>
                )}
              </div>
            );
          })}
        </div>
      ))}
    </div>
  );

  // Rendu principal du composant
  return (
    <div className="solo-container">
      <div className="background-grid">
        <SudokuGrid isBackground={true} />
      </div>

      <DifficultyModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSelectDifficulty={handleSelectDifficulty}
      />

      {!isModalOpen && difficulty && (
        <div className="game-content">
          <SudokuGrid isBackground={false} />
          <div className="actions-button-number">
            <NumberSelector
              onNumberSelect={setSelectedNumber}
              selectedNumber={selectedNumber}
            />

            <div className="game-actions">
              <button className="quit-button" onClick={handleQuit}>Quitter</button>
              <div className="error-counter" style={{ background: '#eee', padding: '5px', borderRadius: '4px', color: 'purple' }}>
                Nombre d'erreurs : {errorCount}
              </div>
              <div className="timer" style={{ background: '#eee', padding: '5px', borderRadius: '4px', marginTop: '5px' }}>
                Temps : {elapsedTime}s
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Solo;
