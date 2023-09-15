package main

import (
	"fmt"
	"io"
	"log"
	"net"
    "bytes"
    "encoding/binary"
    "os"
)

const (
	serverAddr    = ":12345"
	maxPacketSize = 2048
	imageFilePath = "received_image.png"
)

type ScreenServer struct{
    conn *net.Conn
    isReceiving bool
}

func main() {
    ss := &ScreenServer{
        conn: nil,
        isReceiving: false,
    }
    ss.start()
}

func (ss *ScreenServer) start() {
	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server is listening on", serverAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
        ss.conn = &conn
		go ss.readLoop()
	}
}

func (ss *ScreenServer) readLoop() {
    if ss.isReceiving {
        return
    }
    ss.isReceiving = true
    defer func() {
        ss.isReceiving = false
    }()

	buf := new(bytes.Buffer)
	for {
        var size int64
        binary.Read(*ss.conn, binary.BigEndian, &size)
        if size == 0 {
            continue
        }
        log.Printf("Waiting for %d bytes\n", size)
		n, err := io.CopyN(buf, *ss.conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received %d bytes\n", n)

        // write to file
        f, err := os.Create(imageFilePath)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        _, err = f.Write(buf.Bytes())
        if err != nil {
            log.Fatal(err)
        }
	}
}
