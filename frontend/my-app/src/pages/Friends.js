import React, { useState, useEffect, useContext } from "react";
import axios from "axios";
import "./Friends.css";
import { Input } from "../components/Input";
import { AuthContext } from "../contexts/AuthContext";
import { API_BASE_URL } from "../config/api";

const Friends = () => {
    const { user } = useContext(AuthContext);
    const [friendUsername, setFriendUsername] = useState("");

    const [friends, setFriends] = useState([]);
    const [removingIds, setRemovingIds] = useState([]);
    const [showPopup, setShowPopup] = useState(false);


    useEffect(() => {
        if (!user) return;

        axios.get(`${API_BASE_URL}/friends/${user.id}`)
            .then((res) => {
                const friendsData = res.data.friends.map((f) => ({
                    id: f.id,
                    username: f.username,
                    status: "Hors ligne",
                }));
                setFriends(friendsData);
            })
            .catch((err) => console.error("Erreur lors du fetch des amis:", err));
    }, [user]);


    const handleAddFriend = () => {
        if (!friendUsername.trim() || !user) return;

        axios.post(`${API_BASE_URL}/add_friend`, {
            user_id: parseInt(user.id, 10),
            friend_username: friendUsername
        }, {
            headers: { "Content-Type": "application/json" }
        })
            .then((res) => {
                const newFriend = res.data.friend;
                setFriends((prev) => [...prev, newFriend]);
                setFriendUsername("");
                setShowPopup(false);
            })
            .catch(() => alert("Impossible d'ajouter cet utilisateur."));
    };

    const handleRemoveFriend = (id) => {
        if (!user) return;

        setRemovingIds((prev) => [...prev, id]);
        axios.delete(`${API_BASE_URL}/delete_friend/${user.id}/${id}`)
            .then(() => setFriends((prev) => prev.filter((f) => f.id !== id)))
            .catch(() => alert("Impossible de supprimer l'ami pour le moment."))
            .finally(() => setRemovingIds((prev) => prev.filter((rid) => rid !== id)));
    };

    return (
        <div className="friends-container">
            <div className="friends-left">
                <div className="friends-header">
                    <h3>➕ Ajouter un ami</h3>
                    <button className="btn-green btn-left" onClick={() => setShowPopup(true)}>Ajouter</button>

                </div>
            </div>

            <div className="friends-right">
                <h3>👥 Mes Amis</h3>
                <ul>
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
