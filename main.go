package main

import (
	"fmt"
	"p2pstore/core"
	"p2pstore/filehandling"
	"p2pstore/group"
	nodehost "p2pstore/host"
	initnode "p2pstore/init"
	"p2pstore/p2p"
	"time"

	"github.com/libp2p/go-libp2p/core/protocol"
)

const topic = "rex/groupscheme/test"
const service = "rex/service/groupschema/test"

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
	ctx, host, err := p2p.EstablishP2PNode()
	filehandling.NodeHostCtx.Host = host
	filehandling.NodeHostCtx.Ctx = ctx
	if err != nil {
		// In casee of any errors while starting the node, node has to shutdown.
		fmt.Println("[ERROR] - while initializing the node")
		fmt.Println("[WARNING] - node will exit in 10 seconds")
		time.Sleep(10 * time.Second)
	}
	host.SetStreamHandler(protocol.ID(group.GroupJoinRequestProtocol), group.HandleStreamJoinRequest)
	host.SetStreamHandler(protocol.ID(group.GroupJoinReplyProtocol), group.RecieveReply)
	host.SetStreamHandler(protocol.ID(filehandling.FileShareProtocol), filehandling.HandleStreamFileShare)
	host.SetStreamHandler(protocol.ID(filehandling.FileShareMetaDataProtocol), filehandling.HandleStreamFileShareMetaIncomming)
	kadDHT := p2p.HandleDHT(ctx, host)
	pubSub := initnode.SetUpPubSub(ctx, host)
	nodeHost := nodehost.NewP2P(ctx, host, kadDHT, pubSub)
	grp, err := group.JoinGroup(nodeHost, "default", topic)
	group.CurrentGroupRoom = grp
	if err != nil {
		fmt.Println("Error while joining the group")
	}

	go core.HandleInputFromSDI(ctx, host, grp)
	go p2p.DiscoverPeers(ctx, host, kadDHT, service)
	for {
		// to make the app run for ever until closed
	}
}
