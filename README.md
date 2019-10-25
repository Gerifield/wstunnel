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


Working example:

[![asciicast](https://asciinema.org/a/CXIHtBqx2qrnpkR6nreCM3K0C.svg)](https://asciinema.org/a/CXIHtBqx2qrnpkR6nreCM3K0C) 
