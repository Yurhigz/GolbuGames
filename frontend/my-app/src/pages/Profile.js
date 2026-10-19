import React, { useContext, useState } from "react";
import "./Profile.css";
import { AuthContext } from "../contexts/AuthContext";
import { API_BASE_URL, authHeaders } from "../config/api";

const Profile = () => {
    const [newPassword, setNewPassword] = useState("");
    const { user } = useContext(AuthContext);

    const handlePasswordChange = async () => {
        if (!newPassword.trim()) return;

        try {
            const response = await fetch(`${API_BASE_URL}/updateuser`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    ...authHeaders(),
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
                <p>Les statistiques et l’historique de jeu ne sont pas encore disponibles.</p>
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

        </div>
    );
};

export default Profile;
