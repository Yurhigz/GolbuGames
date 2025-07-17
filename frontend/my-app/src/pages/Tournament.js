import React, { useState } from "react";
import { Input } from "../components/Input";
import { Card } from "../components/Card";
import "./Tournament.css";

// Anti-XSS : échappe le HTML
const sanitizeInput = (str) => {
  const temp = document.createElement("div");
  temp.textContent = str;
  return temp.innerHTML;
};

// Autorise lettres, chiffres, espace, tiret, underscore, ponctuation basique
const isValidText = (str) => /^[a-zA-Z0-9\s\-_.,!?]+$/.test(str);

const Tournament = () => {
  const [tournaments, setTournaments] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [error, setError] = useState("");
  const [newTournament, setNewTournament] = useState({
    name: "",
    description: "",
    maxPlayers: 4,
  });

  const handleCreateClick = () => setShowModal(true);

  const handleCloseModal = () => {
    setShowModal(false);
    setNewTournament({ name: "", description: "", maxPlayers: 4 });
    setError("");
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewTournament((prev) => ({
      ...prev,
      [name]: name === "maxPlayers" ? parseInt(value, 10) || 1 : value,
    }));
  };

  const handleCreateTournament = () => {
    const rawName = newTournament.name.trim();
    const rawDescription = newTournament.description.trim();

    const name = sanitizeInput(rawName);
    const description = sanitizeInput(rawDescription);
    const maxPlayers = newTournament.maxPlayers;

    if (!name || !description) {
      setError("Veuillez remplir tous les champs.");
      return;
    }

    if (!isValidText(name)) {
      setError("Le nom du tournoi contient des caractères non autorisés.");
      return;
    }

    if (!isValidText(description)) {
      setError("La description contient des caractères non autorisés.");
      return;
    }

    const nameExists = tournaments.some(
      (t) => t.name.toLowerCase() === name.toLowerCase()
    );

    if (nameExists) {
      setError("Un tournoi avec ce nom existe déjà.");
      return;
    }

    const tournament = {
      name,
      description,
      maxPlayers,
    };

    setTournaments((prev) => [...prev, tournament]);
    handleCloseModal();
  };

  return (
    <div className="tournament-container">
      <h1 className="tournament-title">Tournois</h1>

      <div className="tournament-create">
        <button onClick={handleCreateClick} className="mode-button button">
          Créer un tournoi
        </button>
      </div>

      <div className="tournament-list">
        {tournaments.map((t, index) => (
          <Card key={index} className="tournament-card">
            <div className="tournament-card-content">
              <div>
                <div className="tournament-name">{t.name}</div>
                <div className="tournament-players">
                  {t.maxPlayers} joueurs max
                </div>
              </div>
              <button size="sm" className="mode-button button">
                Rejoindre
              </button>
            </div>
          </Card>
        ))}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>Créer un tournoi</h2>

            {error && <div className="modal-error">{error}</div>}

            <Input
              name="name"
              placeholder="Nom du tournoi"
              value={newTournament.name}
              onChange={handleChange}
              className="modal-input"
              autoComplete="off"
            />
            <Input
              name="description"
              placeholder="Description"
              value={newTournament.description}
              onChange={handleChange}
              className="modal-input"
              autoComplete="off"
            />
            <Input
              name="maxPlayers"
              type="number"
              placeholder="Nombre de joueurs"
              value={newTournament.maxPlayers}
              onChange={handleChange}
              className="modal-input"
              min={1}
              max={128}
            />
            <div className="modal-actions">
              <button onClick={handleCreateTournament} className="mode-button button">
                Créer
              </button>
              <button onClick={handleCloseModal} className="mode-button button">
                Annuler
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Tournament;
