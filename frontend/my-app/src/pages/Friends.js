import React, { useState } from "react";
import "./Friends.css";
import { Input } from "../components/Input";

const Friends = () => {
    const [friendUsername, setFriendUsername] = useState("");
    const [friends, setFriends] = useState([
        { id: 1, username: "Alice", status: "En ligne" },
        { id: 2, username: "Bob", status: "Hors ligne" },
    ]);
    const [removingIds, setRemovingIds] = useState([]);
    const [showPopup, setShowPopup] = useState(false);

    const sentRequests = [
        { id: 1, username: "Lucas" },
        { id: 2, username: "Emma" },
    ];

    const receivedRequests = [
        { id: 3, username: "Noah" },
    ];

    const handleInvite = () => {
        if (friendUsername.trim()) {
            alert(`Invitation envoyée à ${friendUsername}`);
            setFriendUsername("");
        }
    };

    const handleRemoveFriend = (id) => {
        setRemovingIds((prev) => [...prev, id]);
        setTimeout(() => {
            setFriends((prev) => prev.filter((f) => f.id !== id));
            setRemovingIds((prev) => prev.filter((rid) => rid !== id));
        }, 300);
    };

    return (
        <div className="friends-container">
            {/* Partie gauche : Ajouter, demandes */}
            <div className="friends-left">
                <div className="friends-header">
                    <h3>➕ Ajouter un ami</h3>
                    <button className="btn-green btn-left" onClick={() => setShowPopup(true)}>
                        Ajouter
                    </button>
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

            {/* Partie droite : liste des amis */}
            <div className="friends-right">
                <h3>👥 Mes Amis</h3>
                <ul className="friend-list-scroll">
                    {friends.map((friend) => (
                        <li
                            key={friend.id}
                            className={`friend-card ${removingIds.includes(friend.id) ? "removing" : ""}`}
                        >
                            <div className="friend-info">
                                <span className="friend-name">{friend.username}</span>
                                <span
                                    className={`status ${
                                        friend.status === "En ligne" ? "online" : "offline"
                                    }`}
                                >
                                    {friend.status}
                                </span>
                            </div>
                            <button
                                className="remove-friend"
                                onClick={() => handleRemoveFriend(friend.id)}
                                title="Supprimer"
                            >
                                ✖
                            </button>
                        </li>
                    ))}
                    {friends.length === 0 && (
                        <li className="empty-msg">Aucun ami pour le moment</li>
                    )}
                </ul>
            </div>

            {/* Popup */}
            {showPopup && (
                <div className="popup-overlay">
                    <div className="popup">
                        <button className="close-button" onClick={() => setShowPopup(false)}>✖</button>
                        <h3>➕ Ajouter un ami</h3>

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
                            <h4>👥 Liste complète des amis</h4>
                            <ul className="friend-list-scroll" style={{ maxHeight: "200px" }}>
                                {friends.map((friend) => (
                                    <li key={friend.id} className="friend-card small">
                                        <div className="friend-info">
                                            <span className="friend-name">{friend.username}</span>
                                            <span className={`status ${friend.status === "En ligne" ? "online" : "offline"}`}>
                                                {friend.status}
                                            </span>
                                        </div>
                                    </li>
                                ))}
                                {friends.length === 0 && (
                                    <li className="empty-msg">Aucun ami</li>
                                )}
                            </ul>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Friends;
