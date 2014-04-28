package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const servAddr = "localhost:9000"

func validName(input string) bool {
	valid, _ := regexp.Compile(`[^\s,;:]{1,20}`)
	return valid.MatchString(input)
}

func messageServer(conn *net.UDPConn, message string, out chan<- string) {
	if len(message) > 0 {
		_, err := conn.Write([]byte(message))
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}
		println("write to server = ", message)
	}

	reply := make([]byte, 512)
	n, _, err := conn.ReadFromUDP(reply[0:])
	if n == 0 || err != nil {
		println("Read from server failed:", err.Error())
		os.Exit(1)
	}

	out <- string(reply)
}

func handleResponse(conn *net.UDPConn, response string, out chan<- string) {
	fmt.Println(response)
	parts := strings.Split(response, ";")
	if strings.Contains(response, "REGISTERED") {
		go func() { messageServer(conn, "", out) }()
	} else if strings.Contains(response, "ROUND STARTING") {
		go func() { messageServer(conn, fmt.Sprintf("JOIN;%s", parts[1]), out) }()
	} else if strings.Contains(response, "YOUR TURN") {
		go func() { messageServer(conn, fmt.Sprintf("ROLL;%s", parts[1]), out) }()
	} else if strings.Contains(response, "ROLLED") {
		go func() { messageServer(conn, fmt.Sprintf("ANNOUNCE;%s;%s", parts[1], parts[2]), out) }()
	} else {
		go func() { messageServer(conn, "", out) }()
	}
}

func main() {
	var name string
	for validName(name) == false {
		fmt.Print(">>>> ")
		_, err := fmt.Scanf("%s", &name)
		if err != nil {
			println("Reading username failed:", err.Error())
			os.Exit(1)
		}
	}

	msg := fmt.Sprintf("REGISTER;%s", name)
	c := make(chan string)

	udpAddr, err := net.ResolveUDPAddr("udp", servAddr)
	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	go func() { messageServer(conn, msg, c) }()

	for {
		timeout := time.After(30 * time.Second)
		select {
		case answer := <-c:
			handleResponse(conn, answer, c)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
}
