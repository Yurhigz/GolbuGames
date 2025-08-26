import React, { useContext, useState } from "react";
import "./Profile.css";
import { AuthContext } from "../contexts/AuthContext";
import { useRequest } from "../utils/Request";

const Profile = () => {
    const [newPassword, setNewPassword] = useState("");
    const [activeCategory, setActiveCategory] = useState("solo");
    const { user } = useContext(AuthContext);
    const { sendRequest } = useRequest();

    const elo = 1320;

    const scoreHistory = {
        solo: [
            { id: 1, result: "Victoire", date: "2025-08-01", time: "320s" },
            { id: 2, result: "Défaite", date: "2025-07-29", time: "400s" },
        ],
        versus: [
            { id: 1, opponent: "Lucas", result: "Victoire", date: "2025-07-27", time: "290s" },
        ],
        tournois: [
            { id: 1, opponent: "Tournoi S1", result: "3e place", date: "2025-07-15" },
        ],
    };

    const handlePasswordChange = async () => {
        if (!newPassword.trim()) return;

        try {
            const res = await sendRequest("POST", "/updateuser", {
                new_password: newPassword,
            }, true);

            alert(res.message || "Mot de passe modifié avec succès !");
            setNewPassword("");
        } catch (err) {
            alert("Erreur : " + err.message);
        }
    };

    // Format titre pour affichage
    const getTitle = (cat) => {
        if (cat === "versus") return "Multijoueurs";
        return cat.charAt(0).toUpperCase() + cat.slice(1);
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
                    Multijoueurs
                </button>
                <button
                    className={`category-tab ${activeCategory === "tournois" ? "active" : ""}`}
                    onClick={() => setActiveCategory("tournois")}
                >
                    Tournois
                </button>
            </div>

            <div className="history-section">
                <h3>Historique - {getTitle(activeCategory)}</h3>
                <table>
                    <thead>
                        <tr>
                            <th>Date</th>
                            {activeCategory !== "solo" && <th>Adversaire</th>}
                            <th>Résultat</th>
                            {(activeCategory === "solo" || activeCategory === "versus") && <th>Temps</th>}
                        </tr>
                    </thead>
                    <tbody key={activeCategory} className="fade-slide-up">
                        {scoreHistory[activeCategory].map((match) => (
                            <tr key={match.id}>
                                <td>{match.date}</td>
                                {activeCategory !== "solo" && <td>{match.opponent}</td>}
                                <td>{match.result}</td>
                                {(activeCategory === "solo" || activeCategory === "versus") && <td>{match.time}</td>}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Profile;
