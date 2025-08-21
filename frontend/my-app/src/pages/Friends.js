import React, { useState, useEffect } from "react";
import axios from "axios";
import "./Friends.css";
import { Input } from "../components/Input";

const Friends = () => {
    const [friendUsername, setFriendUsername] = useState("");
    const [friends, setFriends] = useState([]);
    const [removingIds, setRemovingIds] = useState([]);
    const [showPopup, setShowPopup] = useState(false);

    // --- Fetch friends au montage ---
    useEffect(() => {
        axios.get("http://localhost:3001/friends/4") // 4 = ID utilisateur courant
            .then((res) => {
                const friendsData = res.data.friends.map((f) => ({
                    id: f.id,
                    username: f.username,
                    status: "Hors ligne",
                }));
                setFriends(friendsData);
            })
            .catch((err) => console.error("Erreur lors du fetch des amis:", err));
    }, []);

    // --- Ajouter un ami ---
    const handleAddFriend = () => {
        if (!friendUsername.trim()) return;

        axios.post(`http://localhost:3001/add_friend/4`, { username: friendUsername })
            .then((res) => {
                // Ajouter à la liste des amis
                setFriends((prev) => [...prev, res.data.friend]);
                setFriendUsername("");
                setShowPopup(false);
            })
            .catch((err) => {
                console.error("Erreur lors de l'ajout de l'ami:", err);
                alert("Impossible d'ajouter cet utilisateur.");
            });
    };

    const handleRemoveFriend = (id) => {
        setRemovingIds((prev) => [...prev, id]);
        axios.delete(`http://localhost:3001/delete_friend/4/${id}`)
            .then(() => setFriends((prev) => prev.filter((f) => f.id !== id)))
            .catch((err) => {
                console.error("Erreur lors de la suppression de l'ami:", err);
                alert("Impossible de supprimer l'ami pour le moment.");
            })
            .finally(() => setRemovingIds((prev) => prev.filter((rid) => rid !== id)));
    };

    return (
        <div className="friends-container">
            {/* Partie gauche : Ajouter, demandes */}
            <div className="friends-left">
                <div className="friends-header">
                    <h3>➕ Ajouter un ami</h3>
                    <button className="btn-green btn-left" onClick={() => setShowPopup(true)}>Ajouter</button>
                </div>
            </div>

            {/* Partie droite : liste des amis */}
            <div className="friends-right">
                <h3>👥 Mes Amis</h3>
                <ul className="friend-list-scroll">
                    {friends.map((friend) => (
                        <li key={friend.id} className={`friend-card ${removingIds.includes(friend.id) ? "removing" : ""}`}>
                            <div className="friend-info">
                                <span className="friend-name">{friend.username}</span>
                                <span className={`status ${friend.status === "En ligne" ? "online" : "offline"}`}>
                  {friend.status}
                </span>
                            </div>
                            <button className="remove-friend" onClick={() => handleRemoveFriend(friend.id)}>✖</button>
                        </li>
                    ))}
                    {friends.length === 0 && <li className="empty-msg">Aucun ami pour le moment</li>}
                </ul>
            </div>

            {/* Popup pour ajouter */}
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
                            <button onClick={handleAddFriend}>Envoyer</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Friends;
