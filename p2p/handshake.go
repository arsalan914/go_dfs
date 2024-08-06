package p2p

// HandshakeFunc...?
type HandshakeFunc func(Peer) error 

func NOPhandshakefunc (Peer) error {
	return nil
}
