import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { Button } from "../components/Button";
import { Card, CardContent } from "../components/Card";
import "./Invite.css";

const Invite = () => {
  const { id } = useParams();
  const [joined, setJoined] = useState(false);

  const handleJoin = () => {
    console.log(`Joueur rejoint la partie avec le code : ${id}`);
    setJoined(true);
    // TODO: Connect to backend WebSocket or REST to join room
  };

  return (
    <div className="invite-container">
      <Card>
        <CardContent className="invite-content">
          <h1 className="invite-title">Invitation à une partie</h1>
          <p className="invite-subtitle">
            Rejoignez une partie avec le code : <span className="room-id">{id}</span>
          </p>
          {!joined ? (
            <Button onClick={handleJoin}>Rejoindre la partie</Button>
          ) : (
            <p className="joined-message">Vous avez rejoint la partie ! En attente de l’autre joueur...</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default Invite;
