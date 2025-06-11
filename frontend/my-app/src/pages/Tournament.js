import React, { useState } from "react";
import { Button } from "../components/Button";
import { Input } from "../components/Input";
import { Card } from "../components/Card";
import "./Tournament.css";

const Tournament = () => {
  const [tournaments, setTournaments] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [newTournament, setNewTournament] = useState({
    name: "",
    description: "",
    maxPlayers: 4,
  });

  const handleCreateClick = () => setShowModal(true);
  const handleCloseModal = () => {
    setShowModal(false);
    setNewTournament({ name: "", description: "", maxPlayers: 4 });
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setNewTournament((prev) => ({ ...prev, [name]: value }));
  };

  const handleCreateTournament = () => {
    if (newTournament.name.trim()) {
      setTournaments((prev) => [...prev, newTournament]);
      handleCloseModal();
    }
  };

  return (
    <div className="tournament-container">
      <h1 className="tournament-title">Tournois</h1>

      <div className="tournament-create">
        <Button onClick={handleCreateClick}>Créer un tournoi</Button>
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
              <Button size="sm">Rejoindre</Button>
            </div>
          </Card>
        ))}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <h2>Créer un tournoi</h2>
            <Input
              name="name"
              placeholder="Nom du tournoi"
              value={newTournament.name}
              onChange={handleChange}
              className="modal-input"
            />
            <Input
              name="description"
              placeholder="Description"
              value={newTournament.description}
              onChange={handleChange}
              className="modal-input"
            />
            <Input
              name="maxPlayers"
              type="number"
              placeholder="Nombre de joueurs"
              value={newTournament.maxPlayers}
              onChange={handleChange}
              className="modal-input"
            />
            <div className="modal-actions">
              <Button onClick={handleCreateTournament}>Créer</Button>
              <Button variant="outline" onClick={handleCloseModal}>
                Annuler
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Tournament;
