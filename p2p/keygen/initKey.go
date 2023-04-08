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

func generateKey() (crypto.PrivKey, error) {
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		fmt.Println("Error while generating a key pair")
		return nil, err
	}
	keyFile, err := os.Create(string(initnode.Nodeinfofolder) + "/" + "init.bin")
	defer keyFile.Close()
	if err != nil {
		fmt.Println("Error while creating cache file to store the key")
		return nil, err
	}
	marshalledPrivKey, err := crypto.MarshalPrivateKey(privKey)
	marshalledPubKey, err := crypto.MarshalPublicKey(pubKey)
	keyDataObj := KeyData{
		PrivKey: marshalledPrivKey,
		PubKey:  marshalledPubKey,
	}
	keyBytes, err := json.Marshal(keyDataObj)
	fmt.Println(keyBytes)
	keyFile.Write(keyBytes)
	fmt.Println(peer.IDFromPrivateKey(privKey))
	return privKey, nil
}
