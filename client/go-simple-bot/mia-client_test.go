package main

import (
	"net"
	"testing"
	"time"
)

func TestValidName(t *testing.T) {
	name := "foo:bar"
	if validName(name) {
		t.Fatalf("Name should't be valid: %s", name)
	}

	name = "foooooooooooooooooooo"
	if validName(name) {
		t.Fatalf("Name should't be valid: %s", name)
	}

	name = "foo bar"
	if validName(name) {
		t.Fatalf("Name should't be valid: %s", name)
	}

	name = "foo,bar"
	if validName(name) {
		t.Fatalf("Name should't be valid: %s", name)
	}

	name = "foo;bar"
	if validName(name) {
		t.Fatalf("Name should't be valid: %s", name)
	}

	name = "foobar5"
	if !validName(name) {
		t.Fatalf("Name should be valid: %s", name)
	}
}

func TestNewConnection(t *testing.T) {
	c := newConnection()
	if c == nil {
		t.Fatal("UDP connection should exists")
	}
}

func TestReadFromServer(t *testing.T) {
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9000})
	if err != nil {
		t.Fatalf("ListenUDP failed: %v", err)
	}
	defer c.Close()

	replies := make(chan string)
	done := make(chan bool)

	cc := newConnection()
	defer cc.Close()

	go func() {
		for {
			timeout := time.After(5 * time.Second)
			select {
			case answer, more := <-replies:
				if !more {
					done <- true
					return
				} else if answer != "foo" && answer != "bar" {
					t.Fatalf("Expected foo or bar, got %s", answer)
				}
			case <-timeout:
				t.Fatal("timed out")
				return
			default:
				done <- true
			}

		}
	}()

	cc.Write([]byte("foo"))
	<-done
	cc.Write([]byte("bar"))
	<-done
}

func TestHandleResponse(t *testing.T) {
	replies := make(chan string)
	done := make(chan bool)
	go func() {
		for {
			timeout := time.After(5 * time.Second)
			select {
			case answer, more := <-replies:
				if !more {
					done <- true
					return
				} else if answer != "JOIN;foo" && answer != "ROLL;foo" && answer != "ANNOUNCE;foo;bar" {
					t.Fatalf("Expected JOIN;foo, got %s", answer)
				}
			case <-timeout:
				t.Fatal("timed out")
				return
			default:
				done <- true
			}

		}
	}()

	handleResponse("ROUND STARTING;foo", replies)
	<-done
	handleResponse("YOUR TURN;foo", replies)
	<-done
	handleResponse("ROLLED;foo;bar", replies)
}

func TestMessageServer(t *testing.T) {
}
