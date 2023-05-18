package group

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (gk *GroupKeyShare) GenerateGroupKey() string {
	return fmt.Sprintf("%s/%s/%s", gk.GroupName, gk.Host.Pretty(), gk.Key)
}

func (gr *GroupRoom) PubLoop() {
	for {
		// fmt.Println("publoop running")
		select {
		case <-gr.psctx.Done():
			fmt.Println("PubLoop Exit")
			return
		case outpacket := <-gr.Outbound:
			// fmt.Println("Outbound message is being processed")

			packetbytes, err := json.Marshal(outpacket)
			if err != nil {
				fmt.Println("[ERROR] - during marhsalling outgoing packet")
				continue
			}

			err = gr.pstopic.Publish(gr.psctx, packetbytes)
			if err != nil {
				fmt.Println("[ERROR] - during publishing the packet to group")
				continue
			}
		}
	}
}

func (gr *GroupRoom) SubLoop() {
	go gr.DisplayMessage()

	for {
		select {
		case <-gr.psctx.Done():
			break
		default:
			message, err := gr.psub.Next(gr.psctx)
			if err != nil {
				close(gr.Inbound)
				return
			}

			inpacket := &Packet{}

			err = json.Unmarshal(message.Data, inpacket)
			if err != nil {
				fmt.Println("Error while unmarshalling the messages")
			}
			if message.ReceivedFrom == gr.SelfId {
				if inpacket.PacketType == "<chat>" {
					continue
				}
			}

			gr.Inbound <- *inpacket
		}
	}
}

func (gr *GroupRoom) PeerList() []peer.ID {
	return gr.pstopic.ListPeers()
}

func (gr *GroupRoom) ExitRoom() {
	defer gr.pscancel()
	endoldsession = true
	gr.State = 0
	gr.psub.Cancel()
	gr.pstopic.Close()
}

func (gr *GroupRoom) UpdateUserName(username string) {
	gr.UserName = username
}

func (gr *GroupRoom) BroadCastHandler() {
	waittime := (rand.Intn(60-20) + 20)
	fmt.Println("Broadcast wait time for this node is", waittime)
	for i := 0; i <= waittime; i++ {
		time.Sleep(1 * time.Second)
		if i == waittime {
			if broadcastrecieved {
				broadcastrecieved = false
				go gr.BroadCastHandler()
				return
			}
			if gr.State == 0 {
				fmt.Println("Ending BroadCast handler for " + gr.GroupName)
				return
			}
			broadCastMessage := &BroadCastMessage{
				PeerId: gr.SelfId,
			}
			brdbytes, err := json.Marshal(broadCastMessage)
			if err != nil {
				fmt.Println("[ERROR] - during marshalling broadcast message")
				continue
			}
			brdpacket := &Packet{
				PacketType: "<brd>",
				Content:    brdbytes,
			}
			gr.Outbound <- *brdpacket
			go gr.BroadCastHandler()
		}

	}
}

func (gr *GroupRoom) BroadCastReplyHandler() {
	brdreplypacket := &BroadCastReplyMessage{
		PeerId: gr.SelfId,
	}
	brdreplybytes, err := json.Marshal(brdreplypacket)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling broadcast reply")
		return
	}
	outpacket := &Packet{
		PacketType: "<brdreply>",
		Content:    brdreplybytes,
	}
	gr.Outbound <- *outpacket

}

func (gr *GroupRoom) DisplayMessage() {
	fmt.Println("Starting DisplayMessage Loop")
	for inpacket := range gr.Inbound {
		switch gr.State {
		case 0:
			fmt.Println("Exiting DisplayMessage Loop for " + gr.GroupName)
			return
		default:
			switch inpacket.PacketType {
			case "<chat>":
				chatmsg := &Chatmessage{}
				err := json.Unmarshal(inpacket.Content, chatmsg)
				if err != nil {
					fmt.Println("[ERROR] - during unmarshalling chat message")
				} else {
					fmt.Println("--------------------------------------------------")
					fmt.Printf("%s: %s\n", chatmsg.SenderName, chatmsg.Message)
					fmt.Println("--------------------------------------------------")
				}
				break
			case "<brd>":
				brdmsg := &BroadCastMessage{}
				err := json.Unmarshal(inpacket.Content, brdmsg)
				if err != nil {
					fmt.Println("[ERROR] - during unmarshalling broadcast message")
				} else {
					broadcastrecieved = true
					fmt.Println("BroadCast recieved")
					fmt.Println("**************************************************")
					fmt.Printf("[ANNOUNCE UR SELF <brd>]: %s\n", brdmsg.PeerId.Pretty())
					fmt.Println("**************************************************")
					ResetPeerTable()
					go gr.BroadCastReplyHandler()
				}
				break
			case "<brdreply>":
				brdreplymsg := &BroadCastReplyMessage{}
				err := json.Unmarshal(inpacket.Content, brdreplymsg)
				brdreplymsg.ADDToPeerTable()
				if err != nil {
					fmt.Println("[ERROR] - during unmarshalling broadcast reply message")
				} else {
					fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")
					fmt.Printf("[I AM ACTIVE]: %s\n", brdreplymsg.PeerId.Pretty())
					fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")
				}
				break
			default:
				fmt.Println("[PANIC] - Unkown packet type recieved")
			}

		}
	}
}
