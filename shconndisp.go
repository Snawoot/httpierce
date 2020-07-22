package main

import (
    "net"
    "sync"
)

type connEntry struct {
    conn net.Conn
    err error
    mux sync.Mutex
    refcount int
}


type SharedConnDispatcher struct {
    address string
    dialer *net.Dialer
    sessions map[string]*connEntry
    sessmux sync.Mutex
}

func NewSharedConnDispatcher(address string, dialer *net.Dialer) *SharedConnDispatcher {
    return &SharedConnDispatcher{
        address: address,
        dialer: dialer,
        sessions: make(map[string]*connEntry),
    }
}

func (e *SharedConnDispatcher) ConnectSession(sess_id string) (net.Conn, error) {
    e.sessmux.Lock()
    entry, ok := e.sessions[sess_id]
    if !ok {
        entry = &connEntry{
            refcount: 1,
        }
        entry.mux.Lock()
        e.sessions[sess_id] = entry
        e.sessmux.Unlock()
        conn, err := e.dialer.Dial("tcp", e.address)
        entry.conn, entry.err = conn, err
        entry.mux.Unlock()
        return conn, err
    } else {
        e.sessmux.Unlock()
        entry.mux.Lock()
        entry.refcount++
        conn, err := entry.conn, entry.err
        entry.mux.Unlock()
        return conn, err
    }
}

func (e *SharedConnDispatcher) DisconnectSession(sess_id string) {
    e.sessmux.Lock()
    entry, ok := e.sessions[sess_id]
    if ok {
        entry.mux.Lock()
        entry.refcount--
        if entry.refcount < 1 {
            delete(e.sessions, sess_id)
            entry.conn.Close()
        }
        e.sessmux.Unlock()
        entry.mux.Unlock()
    } else {
        e.sessmux.Unlock()
    }
}
