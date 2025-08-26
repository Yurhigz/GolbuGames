import React, { useState, useEffect } from "react";
import "./Multi.css";

const WaitingModal = ({ isOpen, opponentFound, countdown }) => {
    if (!isOpen) return null;

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                {!opponentFound ? (
                    <>
                        <h2>🔎 Recherche d'adversaire...</h2>
                        <p>Veuillez patienter, nous cherchons un autre joueur.</p>
                        <div className="loader"></div>
                    </>
                ) : (
                    <>
                        <h2>🎉 Adversaire trouvé !</h2>
                        <p>La partie commence dans :</p>
                        <h3>{`${Math.floor(countdown / 60)
                            .toString()
                            .padStart(2, "0")}:${(countdown % 60)
                            .toString()
                            .padStart(2, "0")}`}</h3>
                    </>
                )}
            </div>
        </div>
    );
};

const Multi = () => {
    const [inviteLink, setInviteLink] = useState("");
    const [copied, setCopied] = useState(false);
    const [ws, setWs] = useState(null);
    const [connected, setConnected] = useState(false);
    const [loading, setLoading] = useState(false);
    const [opponentFound, setOpponentFound] = useState(false);
    const [countdown, setCountdown] = useState(300); // 5 minutes en secondes

    const handleMatchmaking = () => {
        console.log("🔎 Recherche de joueur...");
        setLoading(true);

        return new Promise((resolve, reject) => {
            const socket = new WebSocket("ws://localhost:3001/ws/multi");

            socket.onopen = () => {
                console.log("✅ Connecté au serveur WebSocket");
                setConnected(true);
                setWs(socket);
                resolve(socket);
            };

            socket.onerror = (err) => {
                console.error("❌ Erreur WebSocket :", err.message);
                setLoading(false);
                reject(err);
            };

            socket.onclose = () => {
                console.log("🔌 Connexion fermée");
                setConnected(false);
                setLoading(false);
            };
        });
    };

    const startMatchmaking = async () => {
        try {
            const socket = await handleMatchmaking();

            socket.onmessage = (event) => {
                const message = event.data;
                console.log("Message serveur :", message);

                if (message === "Opponent found... Game starting!") {
                    setOpponentFound(true);
                    setCountdown(20);
                }
            };
        } catch (err) {
            console.error("Impossible de lancer le matchmaking :", err);
        }
    };

    // Gestion du countdown
    useEffect(() => {
        if (!opponentFound) return;

        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer);
                    setLoading(false);
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(timer);
    }, [opponentFound]);

    const handleInvite = () => {
        const link = `${window.location.origin}/multiplayer/invite/${generateRoomId()}`;
        setInviteLink(link);
        setCopied(false);
    };

    const generateRoomId = () => Math.random().toString(36).substr(2, 8);

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

            <WaitingModal
                isOpen={loading}
                opponentFound={opponentFound}
                countdown={countdown}
            />

            <div className="multi-options">
                {renderCard(
                    "Matchmaking",
                    "Rejoignez une partie aléatoire contre un autre joueur en ligne.",
                    <button
                        className="button mode-button"
                        onClick={startMatchmaking}
                        disabled={loading}
                    >
                        {loading ? "Recherche..." : "Trouver un joueur"}
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
                                <input type="text" className="input" value={inviteLink} readOnly />
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
