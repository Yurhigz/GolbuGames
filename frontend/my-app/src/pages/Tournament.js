import React, { useCallback, useState, useEffect } from "react";
import axios from "axios";
import { Input } from "../components/Input";
import { Card } from "../components/Card";
import "./Tournament.css";
import { API_BASE_URL, authHeaders } from "../config/api";


const isValidText = (str) => /^[\p{L}\p{N}\s\-_.,!?]+$/u.test(str);

const Tournament = () => {
    const [tournaments, setTournaments] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [error, setError] = useState("");
    const [loadError, setLoadError] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [newTournament, setNewTournament] = useState({
        name: "",
        description: "",
    });

    // 🔹 Fonction pour récupérer les tournois depuis le backend
    const fetchTournaments = useCallback(async () => {
        try {
            setLoadError("");
            const response = await axios.get(`${API_BASE_URL}/tournaments`);
            if (response.data && response.data.tournaments) {
                setTournaments(response.data.tournaments);
            }
        } catch {
            setLoadError("Impossible de charger les tournois.");
        } finally {
            setIsLoading(false);
        }
    }, []);

    // 🔹 Chargement initial des tournois
    useEffect(() => {
        fetchTournaments();
    }, [fetchTournaments]);

    const handleCreateClick = () => setShowModal(true);

    const handleCloseModal = () => {
        setShowModal(false);
        setNewTournament({ name: "", description: "" });
        setError("");
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setNewTournament((prev) => ({ ...prev, [name]: value }));
    };

    const handleCreateTournament = async () => {
        const rawName = newTournament.name.trim();
        const rawDescription = newTournament.description.trim();

        const name = rawName;
        const description = rawDescription;

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
            (t) => (t.Name || t.name).toLowerCase() === name.toLowerCase()
        );

        if (nameExists) {
            setError("Un tournoi avec ce nom existe déjà.");
            return;
        }

        try {
            await axios.post(`${API_BASE_URL}/add_tournament`, {
                Name: name,
                Description: description,
                StartTime: new Date().toISOString(),
                EndTime: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
            }, {
                headers: authHeaders()
            });

            // 🔹 Recharge la liste complète depuis l’API pour rester synchro
            await fetchTournaments();

            handleCloseModal();
        } catch {
            setError("Impossible de créer le tournoi.");
        }
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
                {isLoading ? (
                    <p>Chargement des tournois…</p>
                ) : loadError ? (
                    <p role="alert">{loadError}</p>
                ) : tournaments.length === 0 ? (
                    <p>Aucun tournoi pour le moment.</p>
                ) : tournaments.map((t) => (
                    <Card key={t.ID || t.id} className="tournament-card">
                        <div className="tournament-card-content">
                            <div>
                                <div className="tournament-name">{t.Name || t.name}</div>
                                <div className="tournament-description">
                                    {t.Description || t.description}
                                </div>
                                {t.StartTime && (
                                    <div className="tournament-dates">
                                        Début : {new Date(t.StartTime).toLocaleDateString()} <br />
                                        Fin : {new Date(t.EndTime).toLocaleDateString()}
                                    </div>
                                )}
                            </div>
                            <button type="button" className="mode-button button" disabled>
                                Participation bientôt disponible
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
                        <div className="modal-actions">
                            <button
                                onClick={handleCreateTournament}
                                className="mode-button button"
                            >
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
