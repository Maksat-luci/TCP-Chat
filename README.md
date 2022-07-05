# TCP-Chat


## About Project
This project consists on recreating the `NetCat` in a `Server-Client` Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open `TCP` connections, send UDP packages, listen on arbitrary `TCP` and UDP ports and many more.

To see more information about NetCat inspect the manual `man nc`.
## Usage
```
$ go run . $port
```
### Example
- `1st` Terminal
```bash
$ go run .
Listening on the port :8080
2021/10/27 19:14:10 Connected 127.0.0.1:39744
2021/10/27 19:14:32 Connected 127.0.0.1:39750
2021/10/27 19:16:10 Connect 127.0.0.1:39750 was left
2021/10/27 19:16:13 Connect 127.0.0.1:39744 was left
^C2021/10/27 19:18:07 Closing Server
2021/10/27 19:18:07 Server Closed
$ 
```
- `2nd` Terminal
```bash
$ nc localhost $yourPort or default port 8989
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Dias1c
[2021-10-27 19:14:10][Dias1c]:Hello? 
[2021-10-27 19:14:20][Dias1c]:!
nrblzn has joined our chat...
[2021-10-27 19:14:32][Dias1c]:!
[2021-10-27 19:14:43][nrblzn]:Hi
[2021-10-27 19:14:43][Dias1c]:Wow, How Are you?
[2021-10-27 19:15:04][Dias1c]:!
[2021-10-27 19:15:31][nrblzn]:I am good, thank you!
[2021-10-27 19:15:31][Dias1c]:!
[2021-10-27 19:15:51][nrblzn]:Lets play Lem-in
[2021-10-27 19:15:51][Dias1c]:HA-ha Go!
[2021-10-27 19:16:03][Dias1c]:!
nrblzn has left our chat...
[2021-10-27 19:16:10][Dias1c]:^C
$ 
```
- `3rd` Terminal
```bash
$ nc localhost 8080
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Amaymon
[2021-10-27 19:14:20][Dias1c]:Hello?
[2021-10-27 19:14:32][Amaymon]:Hi
[2021-10-27 19:14:43][Amaymon]:!
[2021-10-27 19:15:04][Dias1c]:Wow, How Are you?
[2021-10-27 19:15:04][Amaymon]:I am good, thank you!
[2021-10-27 19:15:31][Amaymon]:Lets play Lem-in
[2021-10-27 19:15:51][Amaymon]:!
[2021-10-27 19:16:03][Dias1c]:HA-ha Go!
[2021-10-27 19:16:03][Amaymon]:^C
$ 
```
### Default settings
Default settings on main.go
```go
MaxConnections = 10 // Max connections count
Port = ":8989"      // Default port if user not sets
```
## Build
```bash
#Build project
$ go build -o TCPChat
#Usage
$ ./TCPChat $port
```
## Authors
- [Amaymon](https://github.com/Maksat-luci)