package main

import (
	"fmt"
	"log"
	"net"
	"time"
    "context"
    "sync"
    "github.com/google/uuid"
    "encoding/hex"
    "io"
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
        defer wg.Done()
        forwardUp(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
    }()
    go func() {
        defer wg.Done()
        forwardDown(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
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

    _, err = remoteConn.Write(makeReqBuffer(sess_id, true))
    if err != nil {
        log.Printf("WARN: request write failed: %v", err)
        return
    }
    chunkedWriter := NewChunkedWriter(remoteConn)
    defer chunkedWriter.Close()

    done := make(chan struct{}, 1)
    go func() {
        io.Copy(chunkedWriter, localConn)
        done <- struct{}{}
    }()

    select {
    case <-ctx.Done():
        localConn.SetReadDeadline(epoch)
        <-done
        localConn.SetReadDeadline(zeroTime)
        return
    case <-done:
        return
    }
}

func forwardDown(ctx context.Context, localConn net.Conn, serverAddr string, dialer net.Dialer, sess_id uuid.UUID) {
    remoteConn, err := dialer.DialContext(ctx, "tcp", serverAddr)
    if err != nil {
        log.Printf("WARN: forward downstream connection failed: %v", err)
        return
    }
    defer remoteConn.Close()

    _, err = remoteConn.Write(makeReqBuffer(sess_id, false))
    if err != nil {
        log.Printf("WARN: request write failed: %v", err)
        return
    }
    chunkedReader := NewChunkedReader(remoteConn)

    done := make(chan struct{}, 1)
    go func() {
        io.Copy(localConn, chunkedReader)
        done <- struct{}{}
    }()

    select {
    case <-ctx.Done():
        remoteConn.SetReadDeadline(epoch)
        <-done
        remoteConn.SetReadDeadline(zeroTime)
        return
    case <-done:
        return
    }
}

func makeReqBuffer(sess_id uuid.UUID, upload bool) []byte {
    buf := []byte(fmt.Sprintf("POST #%s# HTTP/1.1\r\n", hex.EncodeToString(sess_id[:])))
    if upload {
        buf = append(buf, header_chunked...)
    }
    buf = append(buf, trailer...)
    return buf
}
