package websocket

//  Gestionnaire central des connexions

// Centre de gestion des connexions
// Gère les rooms (pour le multijoueur)
// Distribue les messages
// Maintient la liste des clients connectés
// structure du hub =>
// Hub
// └── Rooms
//     ├── Room1
//     │   ├── Client1
//     │   └── Client2
//     └── Room2
//         ├── Client3
//         └── Client4

// Il faut considérer un mainhub et des sous hubs dans la mesure où on veut faire cohabiter plusieurs jeux
