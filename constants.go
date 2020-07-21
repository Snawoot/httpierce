package main

import (
    "time"
)

var (
    header_chunked = []byte("Transfer-Encoding: chunked\r\n")
    trailer = []byte("\x20\xd9\x8c\xd9\x8f\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x91\xd9\x91\xd9\x91\xd9\x91\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x8f\xd9\x8f\xd9\x8f\xd9\x92\xd9\x8f\xd9\x8f\r\n\r\n")
    epoch = time.Unix(0, 0)
    zeroTime time.Time
    respDown = []byte("HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n")
    respDownLen = len(respDown)
    respUp = []byte("HTTP/1.1 204 No Content\r\n\r\n")
)
