package group

import "fmt"

func ResetPeerTable() {
	PeerTable = nil
}

func PrintPeerTable() {
	fmt.Println("============================================")
	for _, activePeer := range PeerTable {
		fmt.Println(activePeer.PeerId.Pretty(), activePeer.UserName)
	}
	fmt.Println("============================================")
}
	