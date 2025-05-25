package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func main() {

	// http.HandleFunc("/websocket", websocketHandler)
	// fmt.Printf("Listening on port %v ...", 3005)

	// http.ListenAndServe(":3005", nil)
	// time.Sleep(100 * time.Millisecond)
	// testWebSocketClient()

	// select {}
	b := byte(0b00001011)

	fmt.Println(isFinalFrame(b))
	fmt.Println(getOpcode(b))
	fmt.Println(hasAnyRSVSet(b))
	fmt.Println(buildControlByte(true, false, false, false, 0x1))
}

func isFinalFrame(b byte) bool {
	if b&(1<<7) != 0 {
		return true
	}
	return false
}

func getOpcode(b byte) byte {
	return b & 0b00001111
}

func hasAnyRSVSet(b byte) bool {
	if (b&0b01110000)>>4 > 0 {
		return true
	}
	return false
}

func buildControlByte(fin bool, rsv1, rsv2, rsv3 bool, opcode byte) byte {
	b := byte(0b00000000)
	if fin {
		b |= (1 << 7)
	}
	if rsv1 {
		b |= (1 << 6)
	}
	if rsv2 {
		b |= (1 << 5)

	}
	if rsv3 {
		b |= (1 << 4)
	}
	return b | opcode
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Header.Get("Connection")) != "upgrade" || strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		http.Error(w, "Invalid upgrade request", http.StatusBadRequest)
		return
	}

	secretKey := r.Header.Get("Sec-WebSocket-Key")

	if secretKey == "" {
		http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}

	// un GUID est fixé par le RFC , il est immuable , preuve de conformité lors des handshakes et upgrade vers websocket

	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(secretKey + magicGUID))
	acceptKey := base64.StdEncoding.EncodeToString(hash[:])

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	fmt.Println("WebSocket connection established!")
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		fmt.Printf("Received: %s\n", string(buf[:n]))
		firstByte := buf[0]
		secondByte := buf[1]

		fmt.Printf("BIT FIRSTBYTE : %d", firstByte)
		fmt.Printf("BIT FIRSTBYTE : %d", secondByte)
		// fmt.Printf("Ceci est le premier byte %b", buf[0])
		// fmt.Printf("Ceci est le premier byte en version string %s", string(buf[0]))
	}

}

// func testWebSocketClient() {
// 	conn, err := net.Dial("tcp", "localhost:3005")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	// Envoyer la requête d'upgrade
// 	request := "GET /websocket HTTP/1.1\r\n" +
// 		"Host: localhost:3005\r\n" +
// 		"Upgrade: websocket\r\n" +
// 		"Connection: Upgrade\r\n" +
// 		"Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n" +
// 		"Sec-WebSocket-Version: 13\r\n\r\n"

// 	conn.Write([]byte(request))

// 	// Lire la réponse
// 	buffer := make([]byte, 1024)
// 	n, _ := conn.Read(buffer)
// 	fmt.Println(string(buffer[:n]))
// }
