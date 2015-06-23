package utils

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
)

type zeroFile struct{}

type nullFile struct{}

func (d *nullFile) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d *zeroFile) Read(p []byte) (int, error) {
	return len(p), nil
}

var devNull = &nullFile{}
var devZero = &zeroFile{}

//DevNullConnection take a connection on the receive end, get all data
//and put into an empty reader
func DevNullConnection(conn net.Conn, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	_, err := io.Copy(devNull, conn)
	if err != nil {
		log.Printf("data receive error: %v", err)
		return
	}
	return
}

//SendControlSignal send a message to a TCP address
func SendControlSignal(address string, msg int32) error {
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, msg)
	if err != nil {
		return err
	}
	return nil
}
