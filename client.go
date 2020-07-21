package main

import (
	"fmt"
	"log"
	"net"
	"time"
    "context"
    "sync"
    "github.com/google/uuid"
)

func DoClient(l net.Listener, serverAddr string, timeout time.Duration, vpnMode bool) error {
	dialer := net.Dialer{
		Timeout: timeout,
		Control: GetControlFunc(&TcpConfig{AndroidVPN: vpnMode}),
	}

	for {
		localConn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("l.Accept(): %w", err)
		}

		go func() {
			defer localConn.Close()
            serveConn(localConn, serverAddr, dialer)
		}()
	}
}

func serveConn(localConn net.Conn, serverAddr string, dialer net.Dialer) {
    sess_id := uuid.New()
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        forwardUp(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
        wg.Done()
    }()
    go func() {
        forwardDown(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
        wg.Done()
    }()
    wg.Wait()
}

func forwardUp(ctx context.Context, localConn net.Conn, serverAddr string, dialer net.Dialer, sess_id uuid.UUID) {
    remoteConn, err := dialer.DialContext(ctx, "tcp", serverAddr)
    if err != nil {
        log.Printf("WARN: forward upstream connection failed: %v", err)
        return
    }
    defer remoteConn.Close()
}

func forwardDown(ctx context.Context, localConn net.Conn, serverAddr string, dialer net.Dialer, sess_id uuid.UUID) {
    remoteConn, err := dialer.DialContext(ctx, "tcp", serverAddr)
    if err != nil {
        log.Printf("WARN: forward downstream connection failed: %v", err)
        return
    }
    defer remoteConn.Close()
}
