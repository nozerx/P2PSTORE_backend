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
		// In case of any errors while creation of folders, node cannot function properly, so has to shutdown.
		fmt.Println(err.Error())
		fmt.Println("[WARNING] - node will exit in 10 seconds")
		time.Sleep(10 * time.Second)
		return
	}

	// The following code calls the function EstablishP2PNode() to kickstart the node.
	_, _, err = p2p.EstablishP2PNode()
	if err != nil {
		// In casee of any errors while starting the node, node has to shutdown.
		fmt.Println("[ERROR] - while initializing the node")
		fmt.Println("[WARNING] - node will exit in 10 seconds")
		time.Sleep(10 * time.Second)
	}

	for {
		// to make the app run for ever until closed
	}
}
