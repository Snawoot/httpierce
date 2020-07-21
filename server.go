package main

import (
	"fmt"
	//"log"
	"net"
	"time"
)

func DoServer(l net.Listener, dst string, timeout time.Duration) error {
	for {
		clientConn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("l.Accept(): %w", err)
		}
        defer clientConn.Close()

		go handleClientConn(clientConn, dst, timeout)
	}
}

func handleClientConn(remoteConn net.Conn, dst string, timeout time.Duration) {
}
