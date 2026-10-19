import React, { useState, useEffect, useRef, useCallback } from "react";
import "./Multi.css";
import { WEBSOCKET_BASE_URL } from "../config/api";

const Multi = () => {
  const [status, setStatus] = useState("idle"); // idle | searching | waiting | countdown | playing | finished | left
  const [message, setMessage] = useState("");
  const [winner, setWinner] = useState(null);
  const [grid, setGrid] = useState([]);
  const [givenCells, setGivenCells] = useState([]); // cellules fixes
  const [selectedCell, setSelectedCell] = useState(null);
  const [errors, setErrors] = useState(0);
  const [time, setTime] = useState(0);
  const wsRef = useRef(null);
  const statusRef = useRef("idle");
  const manualCloseRef = useRef(false);

  const updateStatus = (nextStatus) => {
    statusRef.current = nextStatus;
    setStatus(nextStatus);
  };

  // Transforme un tableau plat en matrice 9x9
  const parseBoardArray = (arr) => {
    if (!arr || !Array.isArray(arr) || arr.length !== 81) {
      return Array(9).fill(null).map(() => Array(9).fill(null));
    }
    const matrix = [];
    for (let i = 0; i < 9; i++) {
      matrix.push(
        arr.slice(i * 9, i * 9 + 9).map((num) => (num === 0 ? null : num))
      );
    }
    return matrix;
  };

  // Calcule les cases données (non jouables)
  const computeGivenIndexes = (matrix) => {
    const out = [];
    matrix.forEach((row, r) =>
      row.forEach((val, c) => {
        if (val != null) out.push(r * 9 + c);
      })
    );
    return out;
  };

  const handleNumberSelect = useCallback((number) => {
    if (
      !selectedCell ||
      givenCells.includes(selectedCell.row * 9 + selectedCell.col)
    )
      return;

    const position = selectedCell.row * 9 + selectedCell.col;

    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({ type: "validate_move", position, value: number })
      );
    }
  }, [givenCells, selectedCell]);

  // Timer
  useEffect(() => {
    if (status !== "playing") return;
    const timer = setInterval(() => setTime((t) => t + 1), 1000);
    return () => clearInterval(timer);
  }, [status]);

  // Cleanup WS
  useEffect(() => {
    return () => {
      if (wsRef.current) {
        manualCloseRef.current = true;
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);

  // Contrôle clavier
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (status !== "playing" || !selectedCell) return;

      const { row, col } = selectedCell;

      if (e.key >= "1" && e.key <= "9") {
        handleNumberSelect(parseInt(e.key, 10));
      } else if (e.key === "ArrowUp" && row > 0) {
        setSelectedCell({ row: row - 1, col });
      } else if (e.key === "ArrowDown" && row < 8) {
        setSelectedCell({ row: row + 1, col });
      } else if (e.key === "ArrowLeft" && col > 0) {
        setSelectedCell({ row, col: col - 1 });
      } else if (e.key === "ArrowRight" && col < 8) {
        setSelectedCell({ row, col: col + 1 });
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleNumberSelect, status, selectedCell]);

  const handleMatchmaking = () => {
    const token = localStorage.getItem("token");
    if (!token) {
      updateStatus("idle");
      setMessage("Connectez-vous pour jouer en multijoueur.");
      return;
    }

    if (wsRef.current) {
      const previousSocket = wsRef.current;
      wsRef.current = null;
      previousSocket.close();
    }
    manualCloseRef.current = false;

    const ws = new WebSocket(`${WEBSOCKET_BASE_URL}/ws/multi`, [
      "golbugames",
      `auth.${token}`,
    ]);
    wsRef.current = ws;

    ws.onopen = () => {
      updateStatus("searching");
      setMessage("🔍 Recherche d’un joueur...");
    };

    ws.onmessage = (event) => {
  try {
    const msg = JSON.parse(event.data);
    console.log("MSG RECU:", msg);

    if (msg.type === "system") {
      switch (msg.message) {
        case "waitingOpponent":
          updateStatus("waiting");
          setMessage("En attente d’un autre joueur...");
          break;

        case "opponentFound":
          setMessage("Adversaire trouvé ! Préparation de la partie...");
          break;

        case "gameStarting":
          updateStatus("countdown");
          let count = msg.countdown || 5;
          setMessage(`La partie commence dans ${count}...`);

          const interval = setInterval(() => {
            count -= 1;
            if (count > 0) {
              setMessage(`La partie commence dans ${count}...`);
            } else {
              clearInterval(interval);
              updateStatus("playing");
              setMessage("🎮 Partie commencée !");
              const board = parseBoardArray(msg.grid);
              setGrid(board);
              setGivenCells(computeGivenIndexes(board));
              setErrors(0);
              setTime(0);
            }
          }, 1000);
          break;

        case "gameEnding":
          updateStatus("finished");
          setWinner(msg.winner);
          setMessage(`Partie terminée. Vainqueur: ${msg.winner}`);
          break;

        case "opponentLeft":
          updateStatus("left");
          setMessage("❌ Votre adversaire a quitté la partie.");
          break;

        default:
          console.log("Message système ignoré:", msg);
      }
    }

    if (msg.type === "validate_move") {
      const { position, value, valid } = msg;
      if (!Number.isInteger(position) || position < 0 || position >= 81) return;
      const row = Math.floor(position / 9);
      const col = position % 9;
      if (valid) {
        setGrid((prev) => {
          const newGrid = prev.map((r) => [...r]);
          newGrid[row][col] = value;
          return newGrid;
        });
        setGivenCells((prev) => [...prev, row * 9 + col]);
      } else {
        setErrors((prev) => prev + 1);
      }
    }
  } catch (e) {
    console.error("Erreur parsing message:", e);
  }
};


    ws.onclose = () => {
      if (wsRef.current !== ws || manualCloseRef.current) return;
      if (!["finished", "left"].includes(statusRef.current)) {
        updateStatus("idle");
        setMessage("❌ Déconnecté du serveur.");
      }
    };
  };
        
  const handleQuit = () => {
    if (wsRef.current) {
      manualCloseRef.current = true;
      wsRef.current.close();
      wsRef.current = null;
    }
    updateStatus("idle");
    setMessage("Vous avez quitté la partie.");
  };

  return (
    <div className="multi-container">
      <h1 className="multi-title">Multijoueur Sudoku</h1>

      {status === "idle" && (
        <>
          <button className="button mode-button" onClick={handleMatchmaking}>
            Trouver un joueur
          </button>
          {message && <p className="status-box">{message}</p>}
        </>
      )}

      {["searching", "waiting", "countdown"].includes(status) && (
        <div className="status-box">
          <p>{message}</p>
        </div>
      )}

      {status === "playing" && (
        <div className="flex justify-center items-center">
          {/* Grille */}
<div className="sudoku-grid">
  {grid.map((row, rowIndex) => (
    <div key={rowIndex} className="grid-row">
      {row.map((cell, colIndex) => {
        const pos = rowIndex * 9 + colIndex;
        const isSelected =
          selectedCell?.row === rowIndex && selectedCell?.col === colIndex;
        const isGiven = givenCells.includes(pos);

        return (
          <div
            key={`${rowIndex}-${colIndex}`}
            onClick={() => setSelectedCell({ row: rowIndex, col: colIndex })}
            className={`grid-cell 
              ${isSelected ? "selected-cell" : ""} 
              ${isGiven ? "given-cell" : ""}`}
          >
            <span>{cell ?? ""}</span>
          </div>
        );
      })}
    </div>
  ))}
</div>



          {/* HUD */}
          <div className="ml-6 flex flex-col items-center">
            {/* Pavé numérique */}
            <div className="grid grid-cols-3 gap-2 mb-6">
              {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((n) => (
                <button
                  key={n}
                  onClick={() => handleNumberSelect(n)}
                  className="w-12 h-12 bg-blue-500 text-white font-bold rounded shadow"
                >
                  {n}
                </button>
              ))}
            </div>

            {/* Infos */}
            <div className="text-center mb-4">
              <p>⏱ Temps : {time}s</p>
              <p>❌ Erreurs : {errors}</p>
            </div>

            {/* Quitter */}
            <button
              onClick={handleQuit}
              className="bg-red-500 text-white px-4 py-2 rounded shadow"
            >
              Quitter
            </button>
          </div>
        </div>
      )}

      {status === "finished" && (
        <div className="status-box">
          <p>{message}</p>
          {winner && <p className="winner">🏆 Gagnant : {winner}</p>}
          <button className="button mode-button" onClick={handleQuit}>
            Retour au menu
          </button>
        </div>
      )}

      {status === "left" && (
        <div className="status-box">
          <p>{message}</p>
          <button className="button mode-button" onClick={handleQuit}>
            Retour au menu
          </button>
        </div>
      )}
    </div>
  );
};

export default Multi;
