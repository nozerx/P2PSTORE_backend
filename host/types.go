package host

import (
	"context"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
)

type P2P struct {
	Ctx    context.Context
	Host   host.Host
	KadDHT *dht.IpfsDHT
	PubSub *pubsub.PubSub
}

func NewP2P(ctx context.Context, host host.Host, kad_dht *dht.IpfsDHT, pubSub *pubsub.PubSub) *P2P {
	return &P2P{
		Ctx:    ctx,
		Host:   host,
		KadDHT: kad_dht,
		PubSub: pubSub,
	}
}
