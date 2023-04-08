package main

import (
	"fmt"
	initnode "p2pstore/init"
	"p2pstore/p2p"
	"time"
)

func main() {
	err := initnode.InitFolders()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("[WARNING] - node will exit in 10 seconds")
		time.Sleep(10 * time.Second)
	}
	_, _, err = p2p.EstablishP2PNode()
	if err != nil {
		fmt.Println("[ERROR] - while initializing the node")
	}

	for {
		// to make the app run for ever until closed
	}
}
