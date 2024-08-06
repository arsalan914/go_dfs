package p2p

import "net"

// Message represent any arbitrary data is being sent over
// each transport between two nodes in the network
type RPC struct {
	From 		net.Addr
	Payload 	[]byte
}