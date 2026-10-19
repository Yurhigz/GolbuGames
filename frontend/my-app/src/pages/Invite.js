import React from "react";
import { useParams } from "react-router-dom";
import { Button } from "../components/Button";
import { Card, CardContent } from "../components/Card";
import "./Invite.css";

const Invite = () => {
  const { id } = useParams();

  return (
    <div className="invite-container">
      <Card>
        <CardContent className="invite-content">
          <h1 className="invite-title title">Invitation à une partie</h1>
          <p className="invite-subtitle">
            Rejoignez une partie avec le code : <span className="room-id">{id}</span>
          </p>
          <p role="status">La connexion aux parties sur invitation n’est pas encore disponible.</p>
          <Button type="button" disabled>Rejoindre la partie</Button>
        </CardContent>
      </Card>
    </div>
  );
};

export default Invite;
