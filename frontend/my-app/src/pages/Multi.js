import React, { useState } from "react";
import "./Multi.css";

const Multi = () => {
  const [inviteLink, setInviteLink] = useState("");
  const [copied, setCopied] = useState(false);

  const handleMatchmaking = () => {
    console.log("Recherche de joueur...");
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

  const renderCard = (title, description, content) => (
      <div className="card">
        <div className="card-inner">
          <h2 className="option-title">{title}</h2>
          <p className="option-description">{description}</p>
          {content}
        </div>
      </div>
  );

  return (
      <div className="multi-container">
        <h1 className="multi-title">Multijoueur Sudoku</h1>

        <div className="multi-options">
          {renderCard(
              "Matchmaking",
              "Rejoignez une partie aléatoire contre un autre joueur en ligne.",
              <button className="button mode-button" onClick={handleMatchmaking}>
                Trouver un joueur
              </button>
          )}

          {renderCard(
              "Inviter un ami",
              "Créez un lien d'invitation à partager.",
              <>
                <button className="button mode-button" onClick={handleInvite}>
                  Générer un lien
                </button>

                {inviteLink && (
                    <div className="invite-link-wrapper">
                      <input
                          type="text"
                          className="input"
                          value={inviteLink}
                          readOnly
                      />
                      <button className="button mode-button" onClick={handleCopy}>
                        Copier
                      </button>
                    </div>
                )}

                {copied && <p className="copied-message">Lien copié !</p>}
              </>
          )}
        </div>
      </div>
  );
};

export default Multi;
