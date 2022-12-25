
# net-cat

The project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections

**Chat** port can be changed, default port: *8989*

**Chat** is up to 10 client.


## Usage/Examples
Clone the repository and start the server TCP
```bash
$ go run .
Listening on the port :8989
```
Open another terminal and connect to TCP server

```bash
$ nc localhost 8989
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
[ENTER YOUR NAME]:
```

Enter your name and start chatting


## Authors

- @ddarzox
- @robertt3kuk

# net-cat
