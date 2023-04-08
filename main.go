package main

import (
	"fmt"
	initnode "p2pstore/init"
	"p2pstore/p2p"
	"time"
)

func main() {

	// The following code calls the fuction InitFolders() to initialize all the core folders essential for the node to properly work
	err := initnode.InitFolders()
	if err != nil {
		// In case of any errors while creation of folders, node cannot function properly, so has to shut down
		fmt.Println(err.Error())
		fmt.Println("[WARNING] - node will exit in 10 seconds")
		time.Sleep(10 * time.Second)
		return
	}
	_, _, err = p2p.EstablishP2PNode()
	if err != nil {
		fmt.Println("[ERROR] - while initializing the node")
	}

	for {
		// to make the app run for ever until closed
	}
}
