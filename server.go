package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func main() {
	log.Println("SERVER\n")
	addr := ":9090"
	li, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	defer li.Close()

	log.Printf("Server is running on %s\n", addr)
	log.Println("Ready to accept connections")

	commands := make(chan chan string)

	go storage(commands)

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handle(commands, conn)
	}
}

// sends only single response
func handle(commands chan chan string, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		result := make(chan string)
		commands <- result
		io.WriteString(conn, "command: "+ln+", result: "+<-result+"\n")
	}
	conn.Close()
}

func storage(commands chan chan string) {
	for result := range commands {
		result <- "It works!"
	}
}
