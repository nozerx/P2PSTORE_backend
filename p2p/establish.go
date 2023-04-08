package p2p

import (
	"context"
	"fmt"
	"p2pstore/p2p/keygen"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

// This function starts the P2P node
// It initializes the PeerID, with it creates a new node, as well as gets the background context
// This function returns (context, host (basically the node),Errors if any)
func EstablishP2PNode() (context.Context, host.Host, error) {

	// This code either initializes or retrieves the static PrivateKey which serves as the PeerID for the node
	privKey, err := keygen.RetrieveKey() // whole static behavior of node is handled by this code
	if err != nil {
		fmt.Println("Error while getting the key")
		return nil, nil, err
	}
	identity := libp2p.Identity(privKey) // assings the generated PrivateKey as the PeerID for the node using the identity config option
	fmt.Printf("identity: %v\n", identity)
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), identity) // starts a new node
	if err != nil {
		fmt.Println("Error while starting the p2p node")
		return nil, nil, err
	}
	fmt.Println("Successfull in starting the p2p node")
	ctx := context.Background() // collects the background context, to be used all over the application
	return ctx, host, nil
}
