import React from 'react';
import './DifficultyModal.css';

const DifficultyModal = ({ isOpen, onClose, onSelectDifficulty }) => {
    if (!isOpen) return null;

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>Choisissez la difficulté</h2>
                <div className="difficulty-buttons">
                    <button onClick={() => onSelectDifficulty('Facile')}>Facile</button>
                    <button onClick={() => onSelectDifficulty('Moyen')}>Moyen</button>
                    <button onClick={() => onSelectDifficulty('Difficile')}>Difficile</button>
                </div>
                <button className="close-button" onClick={onClose}>×</button>
            </div>
        </div>
    );
};

export default DifficultyModal;