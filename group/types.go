package group

import (
	"context"

	nodehost "p2pstore/host"

	pbsb "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Packet type can be the following
// <brd> - broadcast
// <brdreply> - broadcastreply
// <chat> - chat message

//Types

type Chatmessage struct {
	Message    string
	SenderID   peer.ID
	SenderName string
}

type Packet struct {
	PacketType string
	Content    []byte
}

type BroadCastMessage struct {
	PeerId peer.ID
}

type BroadCastReplyMessage struct {
	PeerId peer.ID
}

type GroupRoom struct {
	HostP2P   *nodehost.P2P
	GroupName string
	UserName  string
	Inbound   chan Packet
	Outbound  chan Packet
	State     int // 1= active and 0- dead
	SelfId    peer.ID
	psctx     context.Context
	pscancel  context.CancelFunc
	pstopic   *pbsb.Topic
	psub      *pbsb.Subscription
}

type ServicePeer struct {
	Id     int
	PeerId peer.ID
}

type GroupKeyShare struct {
	GroupName string
	Host      peer.ID
	Key       string
}

type JoinRequest struct {
	GroupName string
	Host      peer.ID
	Message   string
}

type JoinRequestReply struct {
	GroupName string
	Host      peer.ID
	To        peer.ID
	Message   string
	Granted   bool
	Key       GroupKeyShare
}

type GroupPeer struct {
	PeerId   peer.ID
	UserName string
}

type MentorInfo struct {
	PeerId  peer.ID
	Host    host.Host
	MentCTX *context.Context
}

// Constants

const defaultroom = "lobby"
const defaultusername = "rex"
const GroupJoinRequestProtocol = "/rex/request"
const GroupJoinReplyProtocol = "rex/reply"
const TestProtocol = "test"

// global vars

var broadcastrecieved bool = false
var resetactivetable bool = false
var PeerTable []GroupPeer = nil

var CurrentGroupRoom *GroupRoom
var endoldsession bool
var Peerlist []ServicePeer
var PauseCLI bool = false
var CurrentGroupShareKey *GroupKeyShare //doesnot hanlde the default key
var MentorInfoObj *MentorInfo
var buff []byte
