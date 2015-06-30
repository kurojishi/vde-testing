package vdetesting

import "strconv"

//Port is a Network Port that Contains the port number
//and the methods to use them
type Port struct {
	port int
}

func (p *Port) String() string {
	return strconv.Itoa(p.port)
}

//Int return the Integer for the Port
func (p *Port) Int() int {
	return p.port
}

//NextPort return you the next port in order
func (p *Port) NextPort(i int) Port {
	next := Port{p.port + i}
	return next
}
