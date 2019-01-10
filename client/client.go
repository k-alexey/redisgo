package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var port = flag.Int("port", 9090, "The port to connect to the server")
var host = flag.String("host", "127.0.0.1", "The host to connect to the server")
var help = flag.Bool("help", false, "Print help")

func init() {
	flag.IntVar(port, "p", 9090, "The port to connect to the server")
	flag.StringVar(host, "h", "127.0.0.1", "The host to connect to the server")
}

func main() {
	flag.Parse()
	if *help == true {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("CLIENT connecting to", fmt.Sprintf("%s:%d", *host, *port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		if text == "\n" {
			continue
		}
		fmt.Fprintf(conn, text)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		split := strings.SplitN(message, "#", 2)
		statusLength, err := strconv.Atoi(split[0])
		if err != nil {
			fmt.Println("Protocol error")
		}
		status := split[1][0:statusLength]
		msg := split[1][statusLength:]

		fmt.Printf("%s %s", status, msg)

	}
}
