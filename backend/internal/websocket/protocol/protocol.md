

# Protocole des messages WebSocket

## GameMessages

| **Frontend**                                                                               | **Backend**                                                                                                               |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
| **isValid**<br>`{ "type": "game_message", "value": <value>, "position": <valuePosition> }` | **validationMessage**<br>`{ "type": "game_message", "value": <value>, "position": <valuePosition>, "valid": true/false }` |

---

## ChatMessages

| **Frontend**                                                                                                     | **Backend**                                                                                                         |
| ---------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **clientMessage**<br>`{ "type": "chat_message", "message": <message>, "sender": <client>, "timestamp": <time> }` | **broadcastMessage**<br>`{ "type": "chat_message", "message": <message>, "sender": <client>, "timestamp": <time> }` |

---

## SystemMessages

| **Frontend**                                                                              | **Backend**                                                                                            |
| ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| *(n/a)                | **waitingOpponent**<br>`{ "type": "system_message", "message": <message>, "code": 2000 }` | **opponentFound**<br>`{ "type": "system_message", "message": <message>, "code": 1002 }`
| *(n/a)*                                                                                | **gameStarting**<br>`{ "type": "system_message", "message": <message>, "grid": <grid>, "code": 1000 }` |
| **gameEnding**<br>`{ "type": "system_message", "winner": <clientid>, "looser": <clientid>,"game_length": <duration>, "code": 1001 }`                                                                                  | **gameEnding**<br>`{ "type": "system_message", "message": <message>, "code": 1001 }`                   |


