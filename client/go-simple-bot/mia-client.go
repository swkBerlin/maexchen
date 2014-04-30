package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"net"
	"time"
)

const serverAddr = "127.0.0.1:9000"

func validName(input string) bool {
	valid, _ := regexp.Compile(`[^\s,;:]{1,20}`)
	return valid.MatchString(input)
}

func newConnection() *net.UDPConn {
	ra, err := net.ResolveUDPAddr("udp4", serverAddr)
	if err != nil {
		log.Fatalf("ResolveUDPAddr failed: %v", err.Error())
	}

	c, err := net.DialUDP("udp4", nil, ra)
	if err != nil {
		log.Fatalf("Dial failed: %v", err.Error())
	}

	return c
}

func readFromServer(c *net.UDPConn, out chan<- string) {
	for {
		reply := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(reply)
		if n == 0 || err != nil {
			log.Fatalf("Read from server failed: %v", err.Error())
		}
		reply = reply[:n]
		log.Printf("Read %q", reply)

		out <- string(reply)
	}
}

func messageServer(c *net.UDPConn, message string) {
	log.Printf("Write to server: %q", message)
	n, err := c.Write([]byte(message))
	if n == 0 || err != nil {
		log.Fatalf("WriteToUDP failed: %v", err.Error())
	}
}

func handleResponse(response string, out chan<- string) {
	log.Printf("Read from server: %q", response)
	parts := strings.Split(response, ";")
	if strings.Contains(response, "REJECTED") {
		log.Fatalf("Registration request rejected.")
	} else if strings.Contains(response, "ROUND STARTING") {
		out <- fmt.Sprintf("JOIN;%s", parts[1])
	} else if strings.Contains(response, "YOUR TURN") {
		out <- fmt.Sprintf("ROLL;%s", parts[1])
	} else if strings.Contains(response, "ROLLED") {
		out <- fmt.Sprintf("ANNOUNCE;%s;%s", parts[1], parts[2])
	}
}

func main() {
	var name string
	for validName(name) == false {
		fmt.Print(">>>> ")
		_, err := fmt.Scanf("%s", &name)
		if err != nil {
			log.Fatalf("Reading username failed:", err.Error())
		}
	}

	conn := newConnection()
	defer conn.Close()

	msg := fmt.Sprintf("REGISTER;%s", name)
	replies := make(chan string)
	messages := make(chan string)

	go readFromServer(conn, replies)
	messageServer(conn, msg)

	for {
		timeout := time.After(30 * time.Second)
		select {
		case answer := <-replies:
			go handleResponse(answer, messages)
		case message := <-messages:
			messageServer(conn, message)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
}
