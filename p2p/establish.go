package p2p

import (
	"context"
	"fmt"
	"p2pstore/p2p/keygen"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func EstablishP2PNode() (context.Context, host.Host, error) {
	privKey, err := keygen.RetrieveKey()
	if err != nil {
		fmt.Println("Error while getting the key")
		return nil, nil, err
	}
	identity := libp2p.Identity(privKey)
	fmt.Printf("identity: %v\n", identity)
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), identity)
	if err != nil {
		fmt.Println("Error while starting the p2p node")
		return nil, nil, err
	}
	fmt.Println("Successfull in starting the p2p node")
	ctx := context.Background()
	return ctx, host, nil
}
