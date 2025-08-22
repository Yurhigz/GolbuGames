// Profile.js
import React, {useContext, useState} from "react";
import "./Profile.css";
import {AuthContext} from "../contexts/AuthContext";

const Profile = () => {
    const [newPassword, setNewPassword] = useState("");
    const [activeCategory, setActiveCategory] = useState("solo");
    const { user } = useContext(AuthContext);

    const elo = 1320;

    const scoreHistory = {
        solo: [
            { id: 1, opponent: "Bot#1", result: "Victoire", date: "2025-08-01" },
            { id: 2, opponent: "Bot#2", result: "Défaite", date: "2025-07-29" },
        ],
        versus: [
            { id: 1, opponent: "Lucas", result: "Victoire", date: "2025-07-27" },
        ],
        tournois: [
            { id: 1, opponent: "Tournoi S1", result: "3e place", date: "2025-07-15" },
        ],
    };

    const handlePasswordChange = async () => {
        if (!newPassword.trim()) return;

        try {
            const response = await fetch("http://localhost:3001/updateuser", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id: parseInt(user.id),
                    new_password: newPassword,
                }),
            });

            if (!response.ok) {
                const errData = await response.text();
                alert("Erreur : " + errData);
                return;
            }

            const data = await response.json();
            alert(data.message);
            setNewPassword("");
        } catch (err) {
            alert("Erreur réseau : " + err.message);
        }
    };

    return (
        <div className="profile-container" style={{ fontSize: '0.9rem' }}>
            <h1 className="profile-title">👤 Mon Profil</h1>

            <div className="elo-section">
                <h2>Elo général : <span>{elo}</span></h2>
            </div>

            <div className="password-section">
                <h3>Changer le mot de passe</h3>
                <input
                    type="password"
                    placeholder="Nouveau mot de passe"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                />
                <button onClick={handlePasswordChange}>Modifier</button>
            </div>

            <div className="score-tabs">
                <button
                    className={`category-tab ${activeCategory === "solo" ? "active" : ""}`}
                    onClick={() => setActiveCategory("solo")}
                >
                    Solo
                </button>
                <button
                    className={`category-tab ${activeCategory === "versus" ? "active" : ""}`}
                    onClick={() => setActiveCategory("versus")}
                >
                    Versus
                </button>
                <button
                    className={`category-tab ${activeCategory === "tournois" ? "active" : ""}`}
                    onClick={() => setActiveCategory("tournois")}
                >
                    Tournois
                </button>
            </div>

            <div className="history-section">
                <h3>Historique - {activeCategory.charAt(0).toUpperCase() + activeCategory.slice(1)}</h3>
                <table>
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Adversaire</th>
                            <th>Résultat</th>
                        </tr>
                    </thead>
                    <tbody key={activeCategory} className="fade-slide-up">
                        {scoreHistory[activeCategory].map((match) => (
                            <tr key={match.id}>
                                <td>{match.date}</td>
                                <td>{match.opponent}</td>
                                <td>{match.result}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Profile;
