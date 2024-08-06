package main

import (
	"fmt"
	"log"

	"github.com/arsalan914/go_dfs/p2p"
)

func main (){

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		HandshakeFunc: p2p.NOPhandshakefunc,
		Decoder: p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func(){
		for{
		msg:= <-tr.Consume()
		fmt.Printf("%+v\n",msg)
		}
	}()
	
	if err := tr.ListenAndAccept();err != nil {
		log.Fatal(err)
	}

	select{}

	fmt.Println("hello worlds")
}