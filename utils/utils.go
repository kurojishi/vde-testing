package utils

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	kb int64 = 1000
	mb int64 = 1000 * kb
	gb int64 = 1000 * mb
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
func DevNullConnection(conn net.Conn, wg *sync.WaitGroup) error {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	_, err := io.Copy(devNull, conn)
	if err != nil {
		return err
	}
	return nil
}

//WaitForControlMessage open a listener on port 8999 to get control messages
func WaitForControlMessage(msg int) error {
	var arrived = false
	local, err := Localv4Addr()
	if err != nil {
		return err
	}
	clistener, err := net.Listen("tcp", local+":8999")
	if err != nil {
		return err
	}
	for !arrived {
		conn, err := clistener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		var buf int32
		binary.Read(conn, binary.LittleEndian, &buf)
		if buf == 2 {
			arrived = true
			clistener.Close()
			log.Printf("control message arrived")
		}
	}
	return nil
}

//SendControlSignal send a message to a TCP address
func SendControlSignal(address string, msg int32) error {
	log.Printf("sending control message to %v", address+":8999")
	conn, err := net.Dial("tcp", address+":8999")
	if err != nil {
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, msg)
	if err != nil {
		return err
	}
	return nil
}

//SendControlSignalUntilOnline repeat SendControlSignal until no error return from Dial
func SendControlSignalUntilOnline(address string, msg int32) {
	for {
		if err := SendControlSignal(address, msg); err != nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			log.Print("control message delivered")
			break
		}
	}
}

//SendData send size data (in megabytes)to the string addr
func SendData(addr string, size int64) error {
	_, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	n, err := io.CopyN(conn, devZero, size*(mb))
	if err != nil {
		return err
	}
	if n != size*mb {
		log.Printf("couldnt send %v Megabytes", float64(n)/float64(mb))
		return nil
	}
	return nil
}

//Localv4Addr get the first local ipv4 address that is not loopback
func Localv4Addr() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4 != nil {
			return ipnet.IP.String(), nil
		}
	}
	err = errors.New("No non local Ip adress found")
	return "", err
}

//Localv6Addr get the first local ipv4 address that is not loopback
func Localv6Addr() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4 == nil {
			return ipnet.IP.String(), nil
		}
	}
	err = errors.New("No non local Ip adress found")
	return "", err
}

//InterfaceAddrv4 Get the ipv4 address of a specific interaface
func InterfaceAddrv4(iface *net.Interface) (string, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4 == nil {
			return ipnet.IP.String(), nil
		}
	}
	err = errors.New("No non local Ip adress found")
	return "", err
}
