// pages/Leaderboard.jsx
import React, { useEffect, useState } from "react";
import { Card, CardContent } from "../components/Card";
import "./Leaderboard.css";

const dummyData = [
  { username: "Alpha", score: 4200 },
  { username: "Beta", score: 3900 },
  { username: "Gamma", score: 3600 },
  { username: "Delta", score: 3400 },
  { username: "Epsilon", score: 3100 },
];

const Leaderboard = () => {
  const [players, setPlayers] = useState([]);

  useEffect(() => {
    // TODO: Replace with real backend fetch
    setPlayers(dummyData);
  }, []);

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
                    <td>{player.score}</td>
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
