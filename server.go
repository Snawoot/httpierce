package main

import (
    "fmt"
    "log"
    "net"
    "time"
    "io"
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

        go handleClientConn(clientConn, connDisp)
    }
}

func handleClientConn(remoteConn net.Conn, dispatcher *SharedConnDispatcher) {
    defer remoteConn.Close()
    remoteAddr := remoteConn.RemoteAddr().String()
    log.Printf("Client %s connected", remoteAddr)
    defer log.Printf("Client %s disconnected", remoteAddr)

    upload, sess_id, err := readClientRequest(remoteConn)
    if err != nil {
        log.Printf("Bad request from client %s: %v", remoteAddr, err)
        return
    }
    log.Printf("Client %s: session ID = %s, upload = %t", remoteAddr, sess_id, upload)

    localConn, err := dispatcher.ConnectSession(sess_id)
    if err != nil {
        log.Printf("Client %s: session dispatch error: %v", remoteAddr, err)
        return
    }

    defer dispatcher.DisconnectSession(sess_id)
    if upload {
        forwardServerUp(remoteConn, localConn)
    } else {
        forwardServerDown(remoteConn, localConn)
    }
}

func forwardServerUp(remoteConn, localConn net.Conn) {
    chunkedReader := NewUnwrappedWire(remoteConn)
    _, err := io.Copy(localConn, chunkedReader)
    log.Printf("Client %s: upload stopped: %v", remoteConn.RemoteAddr().String(), err)
    _, err = remoteConn.Write(respUp)
    if err != nil {
        log.Printf("Client %s: response write error: %v",
                   remoteConn.RemoteAddr().String(),
                   err)
    }
}

func forwardServerDown(remoteConn, localConn net.Conn) {
    _, err := remoteConn.Write(respDown)
    if err != nil {
        log.Printf("Client %s: response write error: %v",
                   remoteConn.RemoteAddr().String(),
                   err)
        return
    }
    chunkedWriter := NewWrappedWire(remoteConn)
    defer chunkedWriter.Close()
    _, err = io.Copy(chunkedWriter, localConn)
    log.Printf("Client %s: download stopped: %v", remoteConn.RemoteAddr().String(), err)
}

func readClientRequest(conn net.Conn) (bool, string, error) {
    verb := make([]byte, 4)
    _, err := io.ReadFull(conn, verb)
    if err != nil {
        return false, "", err
    }
    bytesRead := 4

    var upload bool
    switch string(verb) {
    case "GET ":
        upload = false
    case "POST":
        upload = true
    default:
        return false, "", fmt.Errorf("bad HTTP verb in client request: %#v", verb)
    }

    needRead := reqURILen - 1
    if upload {
        needRead++
    }

    sessionBytes := make([]byte, needRead)
    _, err = io.ReadFull(conn, sessionBytes)
    if err != nil {
        return false, "", err
    }
    bytesRead += needRead

    if upload {
        sessionBytes = sessionBytes[2:]
    } else {
        sessionBytes = sessionBytes[1:]
    }

    var remainderLen int
    if upload {
        remainderLen = verbPostReqLen - bytesRead
    } else {
        remainderLen = verbGetReqLen - bytesRead
    }

    err = discardBytes(conn, int64(remainderLen))
    if err != nil {
        return false, "", err
    }
    return upload, string(sessionBytes), nil
}
