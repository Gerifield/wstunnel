package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

const buffSize = 1024 * 20

var buff = make([]byte, buffSize)

func main() {
	target := flag.String("t", "wss://address.here/socket.io/", "Webocket target (wss://...)")
	secondaryAddr := flag.String("sa", "http://localhost:5000/", "Secondary address")
	verbose := flag.Bool("v", false, "Verbose mode")

	proxyMode := flag.Bool("proxy", false, "Enable proxy mode")
	proxyAddr := flag.String("proxyAddr", ":8080", "Local proxy address (used only in proxy mode)")
	flag.Parse()

	u, err := url.Parse(*target)
	if err != nil {
		fmt.Println("invalid url", err)
		return
	}

	conn, err := wssConnect(u, *verbose)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Proxy mode
	if *proxyMode {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.Write(conn)
			if err != nil {
				fmt.Println(err)
				return
			}

			size, err := conn.Read(buff)
			if err != nil {
				fmt.Println(err)
				return
			}

			hj, ok := w.(http.Hijacker)
			if !ok { // Could not hijack the underlying connection, just dump the buffer
				// In this case we could try to read the buffer with `http.ReadRequest()` and set the header and the content from it
				w.Write(buff[:size])
			}

			pConn, _, err := hj.Hijack()
			if err != nil {
				fmt.Println("local proxy hijack failed", err)
				return
			}
			pConn.Write(buff[:size])

		})
		http.ListenAndServe(*proxyAddr, nil)
		return
	}

	// Simple mode
	size, err := getLocalAddress(conn, *secondaryAddr, *verbose, buff)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(buff[:size]))
}

func getLocalAddress(conn *tls.Conn, address string, verbose bool, response []byte) (int, error) {
	req2, _ := http.NewRequest("GET", address, nil)

	if verbose {
		req2.Write(os.Stdout)
	}

	err := req2.Write(conn)
	if err != nil {
		return 0, errors.Wrap(err, "write http headers failed")
	}

	size, err := conn.Read(response)
	if err != nil {
		return 0, errors.Wrap(err, "read secondary fail")
	}
	return size, nil
}

func wssConnect(u *url.URL, verbose bool) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", u.Host+":443", &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, errors.Wrap(err, "connection failed")
	}

	challengeKey, _ := generateChallengeKey()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header["Upgrade"] = []string{"websocket"}
	req.Header["Connection"] = []string{"Upgrade"}
	req.Header["Sec-WebSocket-Key"] = []string{challengeKey}
	req.Header["Sec-WebSocket-Version"] = []string{"1337"}

	if verbose {
		req.Write(os.Stdout)
	}

	//_, err = conn.Write([]byte(wsscontent))
	err = req.Write(conn)
	if err != nil {
		return conn, errors.Wrap(err, "write websocket headers failed")
	}

	buff := make([]byte, buffSize)
	size, err := conn.Read(buff)
	if err != nil {
		return conn, errors.Wrap(err, "read websocket response failed")
	}
	if verbose {
		fmt.Println(string(buff[:size]))
	}

	return conn, nil
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}
