# Laplace

Laplace is an open-source project to enable screen sharing directly via browser.
Made possible using WebRTC for low latency peer-to-peer connections, and WebSocket implemented in golang for WebRTC signaling.

> Demo video: https://youtu.be/E8cUaPrAlzE

[![Laplace for gaming](./doc/laplace-for-gaming.png)](https://youtu.be/E8cUaPrAlzE)


## Try Demo

For demo, you can visit https://laplace.madeby.monster/

![Website screenshot](./doc/screenshot-website.png)


## Motivation

There are already possible solutions to share your computer screen, e.g. TeamViewer.
But most of them require installations of software or plugins.
What Laplace provides is a simple solution to this problem.
For users wanting to share their screen, all they need to do is to open a website page with their browsers, clicking some buttons, then share some session ID with their peers.
No installation or registration required.

#### Solving the latency problem

This project also serves as a proof-of-concept (PoC) for screen sharing capability directly in browsers based on WebRTC.
Using WebRTC, real-time communication is made possible through peer-to-peer connections.
This proves to be very useful in solving one of the biggest problems is screen streaming: **Latency**. 
The latency represents how long the delay is from the source to transmit to the remote client.
If you notice, this latency problem is usually highlighted by game streaming services, since gameplay relies heavily on the interactivity of inputs and outputs.


#### Low server cost 
This solution also solves the server cost problem, since the expensive operations (encoding and transmission) are done on client browsers.
The server is only needed for serving frontends and for WebRTC signaling.


#### Possible Use Cases

- Game streaming from PC to mobile devices.
- Collaborative work where you need to share your screen with remote coworkers.
- Mirroring presentation slides and demonstrations.


## Installation

Build from source

```bash
$ git clone https://github.com/adamyordan/laplace.git
$ cd laplace && go build -o laplace main.go
$ ./laplace --help
```

OR, pull the pre-built docker image

```bash
$ docker pull adamyordan/laplace
$ docker run adamyordan/laplace ./laplace --help
```


## Program Execution

Executing this project basically serves an HTTP server that will host the frontend and the WebSocket implementation.
Note that you sometimes need to run HTTPs in order for browser to connect to websocket.

```bash
$ ./laplace --help
  -addr string
        Listen address (default "0.0.0.0:443")
  -certFile string
        TLS cert file (default "files/server.crt")
  -keyFile string
        TLS key file (default "files/server.key")
  -tls
        Use TLS (default true)
```

By default, you can run the executable without any argument to listen to TLS port 443.
A self-signed certificate files are provided to ease up development.

```bash
$ ./laplace
2020/03/25 01:01:10 Listening on TLS: 0.0.0.0:443
```

You can then open https://localhost:443/ to view Laplace page.
You may need to add certificate exceptions. In Chrome, you can type `thisisunsafe`.



## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License

[MIT](https://choosealicense.com/licenses/mit/)
