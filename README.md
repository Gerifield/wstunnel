# wstunnel

This is a PoC to solve a websocket based challenge from Hacktivity.
(The original solution was in python, this is a go solution, because I prefer go. :) )

The challenge:
https://challenge.0ang3el.tk/websocket.html


There are more challenges, maybe I'll implement them too!

The logic behind the PoC:
- There's a service with websocket
- There's a proxy in front of the websocket app
- There's a flask app running on the same server

Some proxies don't handle the websocket headers/failure well, so here we send an invalid header to the app which rejects it, but the TLS connection will remain open, so we could send and another (plain http!) connection on the same "tunnel", and boom, we have access to the local service!


## Usage

There are 2 different modes which you could use:

### Simple mode

With simple command line flags you could execute a query through the websocket to a local service.

Example: 

```
$ go run main.go -t "wss://challenge.0ang3el.tk/socket.io/?EIO=3&transport=websocket" -sa "http://localhost:5000/flag" -v
```

The `-v` flag will enable the `verbose` mode which dumps all the requests and responses. Without it you'll only get back the last response (with HTTP headers)!

The `-t` is the websocket target and the `-sa` is the "secondary" local target which you'd like to access. 

Simple mode terminal record:

[![asciicast](https://asciinema.org/a/r0p7l89YIn2t9Ckls5QsLCRAR.svg)](https://asciinema.org/a/r0p7l89YIn2t9Ckls5QsLCRAR)


### Proxy mode

In this mode the program will start a local HTTP server which you could call with your browser or curl or anything and it'll tunnel your query to the target.
(It'll basically dynamically change the `-sa` address and send back the response.)

Example:

```
$ go run main.go -t "wss://challenge.0ang3el.tk/socket.io/?EIO=3&transport=websocket" -v -proxy

[in a different terminal]

$ curl -v 127.0.0.1:8080/flag
```

You could even send multiple requests through the same underlying connection.

Proxy mode terminal record:

[![asciicast](https://asciinema.org/a/QNRYAlqDedYVCSd5uQcwy3pmc.svg)](https://asciinema.org/a/QNRYAlqDedYVCSd5uQcwy3pmc)