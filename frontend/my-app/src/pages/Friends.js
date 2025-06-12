import React, { useState } from "react";
import "./Friends.css";
import { Input } from "../components/Input";

const Friends = () => {
    const [friendUsername, setFriendUsername] = useState("");

    const sentRequests = [
        { id: 1, username: "Lucas" },
        { id: 2, username: "Emma" },
    ];

    const receivedRequests = [
        { id: 3, username: "Noah" },
    ];

    const friends = [
        { id: 1, username: "Alice", status: "En ligne" },
        { id: 2, username: "Bob", status: "Hors ligne" },
    ];

    const handleInvite = () => {
        if (friendUsername.trim()) {
            alert(`Invitation envoyée à ${friendUsername}`);
            setFriendUsername("");
        }
    };

    return (
        <div className="friends-container">
            <div className="friends-left">
                <h3>➕ Inviter un ami</h3>
                <div className="invite-form">
                    <Input
                        type="text"
                        placeholder="Nom d'utilisateur"
                        value={friendUsername}
                        onChange={(e) => setFriendUsername(e.target.value)}
                    />
                    <button onClick={handleInvite}>Envoyer</button>
                </div>

                <div className="section">
                    <h4>📥 Demandes reçues</h4>
                    <ul>
                        {receivedRequests.map((req) => (
                            <li key={req.id}>
                                {req.username}
                                <div className="actions">
                                    <button className="accept">Accepter</button>
                                    <button className="decline">Refuser</button>
                                </div>
                            </li>
                        ))}
                        {receivedRequests.length === 0 && (
                            <li className="empty-msg">Aucune demande reçue</li>
                        )}
                    </ul>
                </div>

                <div className="section">
                    <h4>📤 Demandes envoyées</h4>
                    <ul>
                        {sentRequests.map((req) => (
                            <li key={req.id}>{req.username}</li>
                        ))}
                        {sentRequests.length === 0 && (
                            <li className="empty-msg">Aucune demande envoyée</li>
                        )}
                    </ul>
                </div>
            </div>

            <div className="friends-right">
                <h3>👥 Mes Amis</h3>
                <ul>
                    {friends.map((friend) => (
                        <li key={friend.id} className="friend-card">
                            <div className="friend-info">
                                <span className="friend-name">{friend.username}</span>
                                <span className={`status ${friend.status === "En ligne" ? "online" : "offline"}`}>
                                    {friend.status}
                                </span>
                            </div>
                        </li>
                    ))}
                    {friends.length === 0 && (
                        <li className="empty-msg">Aucun ami pour le moment</li>
                    )}
                </ul>
            </div>
        </div>
    );
};

export default Friends;
