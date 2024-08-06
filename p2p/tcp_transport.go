package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct{
	// conn is the underlying connectin of the peer
	conn 		  net.Conn

	// if we dial and retrieve a conn -> outbound = true
	// if we accept and retrieve a conn -> outbound = false 
	outbound      bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn : conn,
		outbound : outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder		  Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener

	mu            sync.RWMutex //mutex protects peers
	peers         map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.TCPTransportOpts.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop(){
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP: accept error: %s\n",err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handeConn(conn)
	}
}

func (t *TCPTransport) handeConn(conn net.Conn){
	peer := NewTCPPeer(conn, false)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	msg := &Message{}
	// buf := make([]byte, 2000)// bytes.Buffer)
	for {
		// n, err:= conn.Read(buf)
		// if err!=nil{
		// 	fmt.Printf("TCP error: %s\n", err)

		// }

		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}		

		msg.From = conn.RemoteAddr()

		// fmt.Printf("message: %+v\n",buf[:n])
		fmt.Printf("message: %+v\n",msg)

	}
}

// func NewTCPTransport(listenAddr string) Transport{
// 	return &TCPTransport{
// 		listenAddress: listenAddr,
// 	}
// }

// func Test(){
// 	t := NewTCPTransport(":4344").(*TCPTransport)
// 	t.listener.Accept()
// }