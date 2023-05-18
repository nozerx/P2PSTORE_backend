package init

import (
	"context"
	"fmt"

	pbsb "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
)

func SetUpPubSub(ctx context.Context, host host.Host) *pbsb.PubSub {
	pubSubHandler, err := pbsb.NewGossipSub(ctx, host)
	if err != nil {
		fmt.Println("Error during creating a new pubsub handler")
	}
	return pubSubHandler
}
