package keygen

import (
	"encoding/json"
	"fmt"
	"os"
	initnode "p2pstore/init"

	"github.com/libp2p/go-libp2p/core/crypto"
)

// This function returns the file size of init.bin file
func getFileSize() int {
	file, _ := os.Stat(string(initnode.Nodeinfofolder) + "/" + "init.bin")
	return int(file.Size())
}

// This function deals with retrieval or initializing of Private Key
func RetrieveKey() (crypto.PrivKey, error) {
	keyfile, err := os.Open(string(initnode.Nodeinfofolder) + "/" + "init.bin") // Tries to open the init file if it already exits
	if err != nil {
		// If the init.bin file doesn't exist ,then we need to initialize the node, with a node PrivatKey
		fmt.Println("Error while opening the init.txt file")
		keyfile.Close()
		return generateKey() // Deals with initializing the node, with PrivateKey and storing it in the init.bin file
	}
	// Below code deals with reading the init.bin file to retrieve the saved PrivateKey
	keyBytes := make([]byte, getFileSize())
	_, err = keyfile.Read(keyBytes)
	if err != nil {
		fmt.Println("Error while retrieving the key from the init.txt file")
		return nil, err
	}
	fmt.Println("Read ", getFileSize(), " bytes")
	keyData := &KeyData{}
	err = json.Unmarshal(keyBytes, keyData)
	if err != nil {
		fmt.Println("Error while unmarshalling the keyData")
		return nil, err
	}
	privKey, err := crypto.UnmarshalPrivateKey(keyData.PrivKey)
	if err != nil {
		fmt.Println("Error during unmarshalling private key")
	}
	return privKey, nil // return the Retrieved Private Key

}
