package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	print("Error")
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.

	for {
		connection, _ := ln.Accept() //ERROR HERE GO ASK
		conns <- connection

	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
		msgs <- Message{clientid, msg}
	}

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()
	ln, _ := net.Listen("tcp", *portPtr)
	//Create a channel for connections
	conns := make(chan net.Conn, 5)
	//Create a channel for messages
	msgs := make(chan Message, 5)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	//Start accepting connections

	id := 0

	go acceptConns(ln, conns)

	for {
		select {
		case conn := <-conns:

			clients[id] = conn
			go handleClient(clients[id], id, msgs)
			id++

			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:
			fmt.Println(msg)
			for id, conn := range clients {
				if !(msg.sender == id) {
					fmt.Fprintln(conn, msg.message)
				}
			}
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		default:
			print("")
		}

	}
}
