package keygen

import (
	"encoding/json"
	"fmt"
	"os"
	initnode "p2pstore/init"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type KeyData struct {
	PrivKey []byte
	PubKey  []byte
}

// This function generates a new PrivateKey, and saves it in the init.bin file at core/nodeinfo folder
func generateKey() (crypto.PrivKey, error) {
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.RSA, 2048) // generates an RSA key pair
	if err != nil {
		fmt.Println("Error while generating a key pair")
		return nil, err
	}
	keyFile, err := os.Create(string(initnode.Nodeinfofolder) + "/" + "init.bin") // creates a new init.bin file
	defer keyFile.Close()
	if err != nil {
		fmt.Println("Error while creating cache file to store the key")
		return nil, err
	}
	marshalledPrivKey, err := crypto.MarshalPrivateKey(privKey) // marshal the private key
	marshalledPubKey, err := crypto.MarshalPublicKey(pubKey)    // marshall the public key
	keyDataObj := KeyData{
		PrivKey: marshalledPrivKey,
		PubKey:  marshalledPubKey,
	}
	keyBytes, err := json.Marshal(keyDataObj)   // marshall the KeyData object
	keyFile.Write(keyBytes)                     // write it to the file init.bin
	fmt.Println(peer.IDFromPrivateKey(privKey)) // Prints the to be node id to the cmd terminal
	return privKey, nil
}
