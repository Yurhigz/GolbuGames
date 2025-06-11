import React, { useState } from "react";
import { Button } from "../components/Button";
import { Card, CardContent } from "../components/Card";
import { Input } from "../components/Input";
import "./Multi.css";

const Multi = () => {
  const [inviteLink, setInviteLink] = useState("");
  const [copied, setCopied] = useState(false);

  const handleMatchmaking = () => {
    console.log("Recherche de joueur...");
    // À connecter à ton backend de matchmaking
  };

  const handleInvite = () => {
    const link = `${window.location.origin}/multiplayer/invite/${generateRoomId()}`;
    setInviteLink(link);
    setCopied(false);
  };

  const generateRoomId = () => {
    return Math.random().toString(36).substr(2, 8);
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(inviteLink);
    setCopied(true);
  };

  return (
    <div className="multi-container">
      <h1 className="multi-title title">Multijoueur Sudoku</h1>

      <div className="multi-options">
        {/* Matchmaking */}
        <Card>
          <CardContent>
            <h2 className="option-title subtitle">Matchmaking</h2>
            <p className="option-description">
              Rejoignez une partie aléatoire contre un autre joueur en ligne.
            </p>
            <Button onClick={handleMatchmaking}>Trouver un joueur</Button>
          </CardContent>
        </Card>

        {/* Inviter un ami */}
        <Card>
          <CardContent>
            <h2 className="option-title subtitle">Inviter un ami</h2>
            <p className="option-description">
              Créez un lien d'invitation à partager.
            </p>
            <Button onClick={handleInvite}>Générer un lien</Button>

            {inviteLink && (
              <div className="invite-link-wrapper">
                <Input value={inviteLink} readOnly />
                <Button onClick={handleCopy}>Copier</Button>
              </div>
            )}

            {copied && <p className="copied-message">Lien copié !</p>}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Multi;
