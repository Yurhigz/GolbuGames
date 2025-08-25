import React, { useEffect, useState } from "react";
import { Card, CardContent } from "../components/Card";
import { useRequest } from "../utils/Request";
import "./Leaderboard.css";

const Leaderboard = () => {
    const [players, setPlayers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const { sendRequest } = useRequest();

    useEffect(() => {
        const fetchLeaderboard = async () => {
            try {
                const data = await sendRequest("GET", "/leaderboard", null, true);
                setPlayers(data);
            } catch (err) {
                console.error("Erreur lors du fetch leaderboard:", err);
                setError("Impossible de charger le classement.");
            } finally {
                setLoading(false);
            }
        };

        fetchLeaderboard();
    }, [sendRequest]);

    if (loading) return <p>Chargement du classement...</p>;
    if (error) return <p className="error">{error}</p>;

    return (
        <div className="leaderboard-container">
            <h1 className="leaderboard-title title">Classement</h1>
            <div className="leaderboard-wrapper">
                <Card className="leaderboard-card">
                    <CardContent>
                        <table className="leaderboard-table">
                            <thead>
                            <tr>
                                <th>Position</th>
                                <th>Nom</th>
                                <th>Score</th>
                            </tr>
                            </thead>
                            <tbody>
                            {players.map((player, index) => (
                                <tr key={player.username}>
                                    <td>{index + 1}</td>
                                    <td>{player.username}</td>
                                    <td>{player.elo_score}</td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default Leaderboard;
