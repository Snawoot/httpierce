package main

import (
    "fmt"
    //"log"
    "net"
    "time"
)

func DoServer(l net.Listener, dst string, timeout time.Duration) error {
    dialer := net.Dialer{
        Timeout: timeout,
    }
    connDisp := NewSharedConnDispatcher(dst, &dialer)

    for {
        clientConn, err := l.Accept()
        if err != nil {
            return fmt.Errorf("l.Accept(): %w", err)
        }
        defer clientConn.Close()

        go handleClientConn(clientConn, connDisp)
    }
}

func handleClientConn(remoteConn net.Conn, dispatcher *SharedConnDispatcher) {
}
