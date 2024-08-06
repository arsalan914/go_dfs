package p2p

// Peer is an interface that represent the remote node.
type Peer interface{
	Close() error
}

// Transport is anything that handles the communication
// between the nodes in network. THis can be fo the 
// type (TCP, UDP, websockets, etc)
type Transport interface{
	ListenAndAccept() error
	Consume() <-chan RPC
}