package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T){
	ops := TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc : NOPhandshakefunc,
		Decoder:  DefaultDecoder{}, 
	}
	tr:= NewTCPTransport(ops)	

	assert.Equal(t, tr.ListenAddr, ":4000")

	// Server
	assert.Nil(t, tr.ListenAndAccept())

}