package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./go-relay <local_port> <remote_host:remote_port>")
		fmt.Println("Example: ./go-relay 9999 192.168.50.10:3306")
		os.Exit(1)
	}

	localPort := os.Args[1]
	remoteAddr := os.Args[2]
	// Listener
	listener, err := net.Listen("tcp", "0.0.0.0:"+localPort)
	if err != nil {
		fmt.Printf("[!] Failed to bind to port %s: %v\n", localPort, err)
		os.Exit(1)
	}
	fmt.Printf("[*] Go-Relay Listening on 0.0.0.0:%s -> Forwarding to %s\n", localPort, remoteAddr)

	// Accept
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Printf("[-] Failed to accept connection: %v\n", err)
			continue
		}
		fmt.Printf("[+] New connection established from %s\n", client.RemoteAddr())
		// Concurrency
		go handleClient(client, remoteAddr)
	}
}
func handleClient(client net.Conn, remoteAddr string) {
	// Close
	defer client.Close()

	// Dial
	remote, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("[!] Failed to connect to target %s: %v\n", remoteAddr, err)
		return
	}
	// Close
	defer remote.Close()
	// Duplex with buffer of 2 for slower goroutines
	done := make(chan bool, 2)
	// C -> R
	go func() {
		io.Copy(remote, client)
		done <- true
	}()
	// R -> C
	go func() {
		io.Copy(client, remote)
		done <- true
	}()
	<-done
	fmt.Printf("[-] Connection closed for %s\n", client.RemoteAddr())
}
