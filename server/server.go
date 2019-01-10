package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var port = flag.Int("port", 9090, "The port for listening on")
var mode = flag.String("mode", "disk", "The possible storage option")
var help = flag.Bool("help", false, "Print help")

type message struct {
	status  string
	message string
}

type command struct {
	line   string
	result chan message
}

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
	li, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	log.Printf("Server is running on %d port\n", *port)
	log.Println("Ready to accept connections")

	commands := make(chan command)

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

func handle(commands chan command, conn net.Conn) {
	defer log.Println("Disconnected", conn.RemoteAddr())
	log.Println("Connected", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		log.Println(conn.RemoteAddr(), ":", ln)
		command := command{
			line:   ln,
			result: make(chan message)}
		commands <- command
		result := <-command.result
		response := fmt.Sprintf("%d#%s%s\n", len(result.status), result.status, result.message)
		io.WriteString(conn, response)
	}
	conn.Close()
}

func storage(commands chan command) {
	data := make(map[string]string)
	file, err := os.Open("./data.dump")
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&data)
	}
	file.Close()
	log.Printf("Storage: %v", data)

	for command := range commands {
		fields := []string{"", "", ""}
		for i, v := range strings.SplitN(command.line, " ", 3) {
			fields[i] = v
		}
		result := message{"", ""}
		if fields[1] == "" {
			result = message{"ERR", "key required"}
		} else {
			switch strings.ToUpper(fields[0]) {
			case "SET":
				data[fields[1]] = fields[2]
				file, err := os.Create("./data.dump")
				if err == nil {
					encoder := gob.NewEncoder(file)
					encoder.Encode(data)
				}
				file.Close()
				result = message{"OK", ""}
			case "GET":
				if v, ok := data[fields[1]]; ok {
					result = message{"OK", v}
				} else {
					result = message{"ERR", "missing value"}
				}
			case "DEL":
				if _, ok := data[fields[1]]; ok {
					delete(data, fields[1])
					result = message{"OK", ""}
				} else {
					result = message{"ERR", "missing key"}
				}
			default:
				result = message{"ERR", "unknown command"}
			}
		}
		command.result <- result
	}
}
