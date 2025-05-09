import React, { useState } from 'react';
import DifficultyModal from '../components/DifficultyModal';
import NumberSelector from '../components/NumberSelector';
import './Solo.css';

const Solo = () => {
    const [isModalOpen, setIsModalOpen] = useState(true);
    const [difficulty, setDifficulty] = useState(null);
    const [selectedNumber, setSelectedNumber] = useState(null);
    const [grid, setGrid] = useState(Array(9).fill(null).map(() => Array(9).fill(null)));

    const handleSelectDifficulty = (selectedDifficulty) => {
        setDifficulty(selectedDifficulty);
        setIsModalOpen(false);
    };

    const SudokuGrid = ({ isBackground }) => (
        <div className={`sudoku-grid ${isBackground ? 'background' : 'active'}`}>
            {[...Array(9)].map((_, rowIndex) => (
                <div key={rowIndex} className="grid-row">
                    {[...Array(9)].map((_, colIndex) => (
                        <div 
                            key={`${rowIndex}-${colIndex}`} 
                            className="grid-cell"
                            onClick={!isBackground ? () => handleCellClick(rowIndex, colIndex) : undefined}
                        >
                            {!isBackground && grid[rowIndex][colIndex]}
                        </div>
                    ))}
                </div>
            ))}
        </div>
    );

    const handleCellClick = (row, col) => {
        if (selectedNumber !== null) {
            const newGrid = [...grid];
            newGrid[row][col] = selectedNumber;
            setGrid(newGrid);
        }
    };

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
                    <NumberSelector 
                        onNumberSelect={setSelectedNumber}
                        selectedNumber={selectedNumber}
                    />
                </div>
            )}
        </div>
    );
};

export default Solo;