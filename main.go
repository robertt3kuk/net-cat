package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	name string
	addr string
}

var (
	clients       = make(map[string]Client)
	leaving       = make(chan message)
	messages      = make(chan message)
	joining       = make(chan message)
	ClientsInfo   = []Client{}
	clientsNumber = 0
)

type message struct {
	text    string
	address string
	name    string
	time    string
}

func main() {
	greetText := greet()
	arg := os.Args[1:]
	portString := ""
	if len(arg) == 0 {
		portString = "8989"
	}
	if len(arg) == 1 {
		if _, err := strconv.Atoi(arg[0]); err != nil {
			fmt.Printf("%q doesn't look like a number.\n", arg[0])
			return
		}
		portString = arg[0]
	}
	if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	listen, err := net.Listen("tcp", "localhost:"+portString)
	fmt.Println("Listening on the port :" + portString)

	if err != nil {
		log.Fatal(err)
	}
	history := createHistoryFile()
	go broadcaster(history)
	defer history.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		go handle(conn, greetText, history, &clientsNumber)
	}
}

func greet() string {
	file, err := os.Open("greet.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	return string(b)
}

func createHistoryFile() *os.File {
	historyFile, err := os.Create("history.txt")
	if err != nil {
		log.Fatal(err)
	}
	return historyFile
}

func handle(conn net.Conn, greet string, history *os.File, counter *int) {
	var client Client
	if len(clients) > 9 {
		fmt.Fprint(conn, "Room full of people, please wait someone to leave")
		conn.Close()
		return
	}
	*counter++
	fmt.Fprint(conn, greet)
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	client.conn = conn
	client.name = strings.TrimSpace(name)
	client.addr = conn.RemoteAddr().String()
	ClientsInfo = append(ClientsInfo, client)
	clients[conn.RemoteAddr().String()] = client
	fi, err := history.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() != 0 {
		text, err := os.ReadFile("history.txt")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(conn, string(text))

	}

	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05")
	joining <- newMessage("has joined our chat...\n", client, timeString)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		timeNow := time.Now()
		timeString := timeNow.Format("2006-01-02 15:04:05")
		if len(strings.TrimSpace(input.Text())) != 0 && input.Text()[0] != 27 {
			messages <- newMessage(input.Text(), client, timeString)
		} else {
			continue
		}
	}

	// Delete client form map
	delete(clients, conn.RemoteAddr().String())

	leaving <- newMessage("has left our chat...\n", client, timeString)
	*counter--
	conn.Close() // NOTE: ignoring network errors
}

func newMessage(msg string, cl Client, time string) message {
	addr := cl.conn.RemoteAddr().String()
	name := cl.name
	return message{
		text:    msg,
		address: addr,
		name:    name,
		time:    time,
	}
}

func broadcaster(history *os.File) {
	for {
		select {
		case msg := <-messages:
			for _, cl := range clients {
				if msg.address != cl.addr {
					fmt.Fprintf(cl.conn, "\n[%s][%s]:%s", msg.time, msg.name, strings.TrimSpace(msg.text)) // NOTE: ignoring network errors
					timeNow := time.Now()
					timeString := timeNow.Format("2006-01-02 15:04:05")
					fmt.Fprintf(cl.conn, "\n[%s][%s]:", timeString, cl.name) // NOTE: ignoring network errors
				} else {
					timeNow := time.Now()
					timeString := timeNow.Format("2006-01-02 15:04:05")
					fmt.Fprintf(cl.conn, "[%s][%s]:", timeString, cl.name)
				}
			}
			history.WriteString(fmt.Sprintf("[%s][%s]:%s\n", msg.time, msg.name, msg.text))
			fmt.Printf("[%s][%s]:%s\n", msg.time, msg.name, msg.text)
		case msg := <-joining:
			for _, cl := range clients {
				if msg.address != cl.addr {
					fmt.Fprintf(cl.conn, "\n%s %s[%s][%s]:", msg.name, msg.text, msg.time, cl.name) // NOTE: ignoring network errors
				} else {
					fmt.Fprintf(cl.conn, "[%s][%s]:", msg.time, msg.name) // NOTE: ignoring network errors
				}
			}
			fmt.Print(msg.name + " " + msg.text)

		case msg := <-leaving:
			for _, cl := range clients {
				if msg.address != cl.addr {
					fmt.Fprintf(cl.conn, "\n%s %s[%s][%s]:", msg.name, msg.text, msg.time, cl.name) // NOTE: ignoring network errors
				} else {
					fmt.Fprintf(cl.conn, "%s %s[%s][%s]:", msg.name, msg.text, msg.time, msg.name) // NOTE: ignoring network errors
				}
			}
			fmt.Print(msg.name + " " + msg.text)
		}
	}
}
