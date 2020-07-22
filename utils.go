package main

import (
    "io"
    "io/ioutil"
    "encoding/base64"
)

func discardBytes(r io.Reader, n int64) error {
    _, err := io.CopyN(ioutil.Discard, r, n)
    return err
}

type WrappedWire struct {
    chunked, radix io.WriteCloser
}

func NewWrappedWire(w io.Writer) *WrappedWire {
    chunked := NewChunkedWriter(w)
    radix := base64.NewEncoder(base64.URLEncoding, chunked)
    return &WrappedWire{chunked, radix}
}

func (ww *WrappedWire) Write(p []byte) (n int, err error) {
    return ww.radix.Write(p)
}

func (ww *WrappedWire) Close() error {
    var err error
    err1 := ww.radix.Close()
    if err1 != nil {
        err = err1
    }
    err2 := ww.chunked.Close()
    if err2 != nil && err != nil {
        err = err2
    }
    return err
}

type UnwrappedWire struct {
    chunked, radix io.Reader
}

func NewUnwrappedWire(r io.Reader) *UnwrappedWire {
    chunked := NewChunkedReader(r)
    radix := base64.NewDecoder(base64.URLEncoding, chunked)
    return &UnwrappedWire{chunked, radix}
}

func (uw *UnwrappedWire) Read(p []byte) (n int, err error) {
    return uw.radix.Read(p)
}
