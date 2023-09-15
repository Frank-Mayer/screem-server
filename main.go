package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	serverAddr    = ":0"
	maxPacketSize = 2048
)

type ScreemServer struct {
	hostConn  net.Conn
	guestConn net.Conn
	ack       chan bool
}

func main() {
	ss := &ScreemServer{
		hostConn:  nil,
		guestConn: nil,
		ack:       make(chan bool),
	}
	ss.start()
	select {}
}

func (self *ScreemServer) start() {
	go self.waitForHost()
	go self.waitForGuest()
}

func (self *ScreemServer) waitForHost() {
	hostLn, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Connect your host to", hostLn.Addr().String())

	for {
		conn, err := hostLn.Accept()
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("Host connected")
		self.hostConn = conn
		go self.readLoop()
	}
}

func (self *ScreemServer) waitForGuest() {
	guestLn, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Connect your guest to", guestLn.Addr().String())

	for {
		conn, err := guestLn.Accept()
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("Guest connected")
		self.guestConn = conn
		go self.ackLoop()
	}
}

func (self *ScreemServer) readLoop() {
	var size int64
	for {
		binary.Read(self.hostConn, binary.BigEndian, &size)
		if size == 0 {
			continue
		}

		if self.guestConn == nil {
			// discard the data
			io.CopyN(io.Discard, self.hostConn, size)
			// send OK to host
			binary.Write(self.hostConn, binary.BigEndian, true)
			continue
		}

		// send data size to guest
		binary.Write(self.guestConn, binary.BigEndian, size)

		// send data to guest
		_, err := io.CopyN(self.guestConn, self.hostConn, size)
		if err != nil {
			log.Panicln(err)
		}

		// await ack from guest
		<-self.ack

		// send OK to host
		binary.Write(self.hostConn, binary.BigEndian, true)
	}
}

func (self *ScreemServer) ackLoop() {
	for {
		var ack bool
		binary.Read(self.guestConn, binary.BigEndian, &ack)
		if ack {
			self.ack <- true
		}
	}
}
