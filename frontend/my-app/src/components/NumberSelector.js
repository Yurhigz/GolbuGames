import React from 'react';
import './NumberSelector.css';

const NumberSelector = ({ onNumberSelect, selectedNumber }) => {
    const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9];

    return (
        <div className="number-selector">
            {numbers.map((number) => (
                <button
                    key={number}
                    className={`number-button ${selectedNumber === number ? 'selected' : ''}`}
                    onClick={() => onNumberSelect(number)}
                >
                    {number}
                </button>
            ))}
        </div>
    );
};

export default NumberSelector;