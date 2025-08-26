import React, { useState, useEffect, useContext } from "react";
import "./Friends.css";
import { Input } from "../components/Input";
import { AuthContext } from "../contexts/AuthContext";
import { useRequest } from "../utils/Request";

const Friends = () => {
    const { user } = useContext(AuthContext);
    const { sendRequest } = useRequest();
    const [friendUsername, setFriendUsername] = useState("");
    const [friends, setFriends] = useState([]);
    const [removingIds, setRemovingIds] = useState([]);
    const [showPopup, setShowPopup] = useState(false);

    useEffect(() => {
        if (!user) return;

        const fetchFriends = async () => {
            try {
                const res = await sendRequest("GET", `/friends`, {}, true);
                const friendsData = res.data.friends.map(f => ({
                    id: f.id,
                    username: f.username,
                    status: "Hors ligne",
                }));
                setFriends(friendsData);
            } catch (err) {
                console.error("Erreur lors du fetch des amis:", err);
            }
        };

        fetchFriends();
    }, [user, sendRequest]);



    const handleAddFriend = async () => {
        if (!friendUsername.trim() || !user) return;

        try {
            const res = await sendRequest("POST", "/add_friend", {
                friend_username: friendUsername
            }, true);

            const newFriend = {
                id: res.friend.id,
                username: res.friend.username,
                status: "Hors ligne"
            };

            setFriends((prev) => [...prev, newFriend]);
            setFriendUsername("");
            setShowPopup(false);
        } catch (err) {
            console.error("Impossible d'ajouter cet utilisateur:", err);
            alert("Impossible d'ajouter cet utilisateur.");
        }
    };

    const handleRemoveFriend = async (id) => {
        if (!user) return;

        setRemovingIds((prev) => [...prev, id]);
        try {
            await sendRequest("DELETE", `/delete_friend/${id}`, {}, true);
            setFriends((prev) => prev.filter((f) => f.id !== id));
        } catch (err) {
            console.error("Impossible de supprimer l'ami:", err);
            alert("Impossible de supprimer l'ami pour le moment.");
        } finally {
            setRemovingIds((prev) => prev.filter((rid) => rid !== id));
        }
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
                <ul className="friend-list-scroll">
                    {friends.length === 0 && <li className="empty-msg">Aucun ami pour le moment</li>}
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
