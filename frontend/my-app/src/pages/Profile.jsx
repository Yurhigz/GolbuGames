import React, { useState } from "react";
import "./Profile.css";

const Profile = () => {
    const [newPassword, setNewPassword] = useState("");
    const [elo] = useState(1320); // exemple Elo statique
    const gamesHistory = [
        { id: 1, opponent: "Alpha", result: "Victoire", date: "2025-08-01" },
        { id: 2, opponent: "Beta", result: "Défaite", date: "2025-07-30" },
        { id: 3, opponent: "Gamma", result: "Victoire", date: "2025-07-28" },
    ];

    const handlePasswordChange = () => {
        if (newPassword.trim()) {
            alert("Mot de passe modifié !");
            setNewPassword("");
        }
    };

    return (
        <div className="profile-container">
            <h1 className="profile-title">👤 Mon Profil</h1>

            <div className="elo-section">
                <h2>Elo : <span>{elo}</span></h2>
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

            <div className="history-section">
                <h3>Historique des parties</h3>
                <table>
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Adversaire</th>
                            <th>Résultat</th>
                        </tr>
                    </thead>
                    <tbody>
                        {gamesHistory.map(game => (
                            <tr key={game.id}>
                                <td>{game.date}</td>
                                <td>{game.opponent}</td>
                                <td>{game.result}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Profile;
