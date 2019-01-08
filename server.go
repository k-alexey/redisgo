package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var port = flag.Int("port", 9090, "Server port")
var mode = flag.String("mode", "disk", "Storage mode")
var help = flag.Bool("help", false, "Print help")

func init() {
	flag.IntVar(port, "p", 9090, "Server port")
	flag.StringVar(mode, "m", "127.0.0.1", "Server host")
}

func main() {
	flag.Parse()
	if *help == true {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Println("SERVER\n")
	li, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	log.Printf("Server is running on %d port\n", *port)
	log.Println("Ready to accept connections")

	commands := make(chan chan string)

	go storage(commands)

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handle(commands, conn)
	}
}

// sends only single response
func handle(commands chan chan string, conn net.Conn) {
	defer log.Println("Disconnected", conn.RemoteAddr())
	log.Println("Connected", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		log.Println(conn.RemoteAddr(), ":", ln)
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
