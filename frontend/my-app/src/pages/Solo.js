import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';

const Solo = () => {
  const [isModalOpen, setIsModalOpen] = useState(true);
  const [difficulty, setDifficulty] = useState(null);
  const [selectedNumber, setSelectedNumber] = useState(null);
  const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));
  const [selectedCell, setSelectedCell] = useState({ row: 0, col: 0 });
  const [result, setResult] = useState(null);
  const navigate = useNavigate();

  const handleSelectDifficulty = (selectedDifficulty) => {
    setDifficulty(selectedDifficulty);
    setIsModalOpen(false);
  };

  const handleCellClick = (row, col) => {
    setSelectedCell({ row, col });
    if (selectedNumber !== null) {
      fillCell(row, col, selectedNumber);
    }
  };

  const fillCell = (row, col, number) => {
    const newGrid = grid.map(row => [...row]);
    newGrid[row][col] = number;
    setGrid(newGrid);
  };

  const isGridComplete = () => {
    return grid.every(row => row.every(cell => cell !== null && cell !== ''));
  };

  const handleSubmit = () => {
    const correct = Math.random() > 0.5;
    setResult(correct ? '✅ Grille correcte !' : '❌ Erreur dans la grille.');
  };

  const handleQuit = () => {
    navigate('/');
  };

  // ⌨️ Gestion clavier
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (isModalOpen) return;

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

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectedCell, isModalOpen, grid]);

  const SudokuGrid = ({ isBackground }) => (
    <div className={`sudoku-grid ${isBackground ? 'background' : 'active'}`}>
      {[...Array(9)].map((_, rowIndex) => (
        <div key={rowIndex} className="grid-row">
          {[...Array(9)].map((_, colIndex) => {
            const isSelected = selectedCell.row === rowIndex && selectedCell.col === colIndex;
            return (
              <div
                key={`${rowIndex}-${colIndex}`}
                className={`grid-cell ${isSelected && !isBackground ? 'selected-cell' : ''}`}
                onClick={!isBackground ? () => handleCellClick(rowIndex, colIndex) : undefined}
              >
                {!isBackground && grid[rowIndex][colIndex]}
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
              <button
                className="submit-button"
                onClick={handleSubmit}
                disabled={!isGridComplete()}
              >
                Soumettre
              </button>
            </div>
          </div>

          {result && <div className="game-result">{result}</div>}
        </div>
      )}
    </div>
  );
};

export default Solo;
