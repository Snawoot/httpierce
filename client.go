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
    log.Printf("Client %s connected", localConn.RemoteAddr().String())
    sess_id := uuid.New()
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        forwardClientUp(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
    }()
    go func() {
        defer wg.Done()
        forwardClientDown(ctx, localConn, serverAddr, dialer, sess_id)
        cancel()
    }()
    wg.Wait()
    log.Printf("Client %s disconnected", localConn.RemoteAddr().String())
}

func forwardClientUp(ctx context.Context, localConn net.Conn, serverAddr string, dialer net.Dialer, sess_id uuid.UUID) {
    remoteConn, err := dialer.DialContext(ctx, "tcp", serverAddr)
    if err != nil {
        select {
        case <-ctx.Done():
        default:
            log.Printf("WARN: forward upstream connection failed: %v", err)
        }
        return
    }

    done := make(chan struct{})
    go func() {
        defer func () {
            done <- struct{}{}
        }()
        _, err = remoteConn.Write(makeReqBuffer(sess_id, true))
        if err != nil {
            select {
            case <-ctx.Done():
            default:
                log.Printf("WARN: request write failed: %v", err)
            }
            return
        }

        chunkedWriter := NewChunkedWriter(remoteConn)
        defer chunkedWriter.Close()
        io.Copy(chunkedWriter, localConn)
    }()

    select {
    case <-ctx.Done():
        localConn.SetReadDeadline(epoch)
        remoteConn.Close()
        <-done
        localConn.SetReadDeadline(zeroTime)
        return
    case <-done:
        remoteConn.Close()
        return
    }
}

func forwardClientDown(ctx context.Context, localConn net.Conn, serverAddr string, dialer net.Dialer, sess_id uuid.UUID) {
    remoteConn, err := dialer.DialContext(ctx, "tcp", serverAddr)
    if err != nil {
        select {
        case <-ctx.Done():
        default:
            log.Printf("WARN: forward downstream connection failed: %v", err)
        }
        return
    }

    done := make(chan struct{})
    go func() {
        defer func () {
            done <- struct{}{}
        }()
        _, err = remoteConn.Write(makeReqBuffer(sess_id, false))
        if err != nil {
            select {
            case <-ctx.Done():
            default:
                log.Printf("WARN: request write failed: %v", err)
            }
            return
        }

        err := discardBytes(remoteConn, int64(respDownLen))
        if err != nil {
            select {
            case <-ctx.Done():
            default:
                log.Printf("WARN: response read failed: %v", err)
            }
            return
        }
        chunkedReader := NewChunkedReader(remoteConn)

        io.Copy(localConn, chunkedReader)
    }()

    select {
    case <-ctx.Done():
        remoteConn.Close()
        localConn.SetWriteDeadline(epoch)
        <-done
        localConn.SetWriteDeadline(zeroTime)
        return
    case <-done:
        remoteConn.Close()
        return
    }
}

func makeReqBuffer(sess_id uuid.UUID, upload bool) []byte {
    var method string
    if upload {
        method = "POST"
    } else {
        method = "GET"
    }
    buf := []byte(fmt.Sprintf("%s #%s# HTTP/1.1\r\n", method, hex.EncodeToString(sess_id[:])))
    if upload {
        buf = append(buf, header_chunked...)
    }
    buf = append(buf, trailer...)
    return buf
}
