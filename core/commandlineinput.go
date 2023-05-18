package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"p2pstore/group"
)

func HandleInputFromSDI(ctx context.Context, host host.Host, gr *group.GroupRoom) {
	reader := bufio.NewReader(os.Stdin)
	for {
		if group.PauseCLI == false {
			gr = group.CurrentGroupRoom
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error during reading input from the stream")
				continue
			}
			if input[:5] == "<cmd>" {
				fmt.Println("These are the available commands")
				fmt.Println("1.Change Group\n2.Change UserName\n3.List Group Peers\n4.List service peers\n5.Join Group\n7.Print PEER-TABLE")
				var choice int
				fmt.Scanln(&choice)
				switch choice {
				case 1:
					group.MentorInfoObj = &group.MentorInfo{
						PeerId:  gr.SelfId,
						Host:    gr.HostP2P.Host,
						MentCTX: &ctx,
					}
					fmt.Println("Enter the new group name to join it")
					groupName, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error during reading input from the SDI")
						continue
					}
					escapeSeqLen := 0
					if runtime.GOOS == "windows" {
						escapeSeqLen = 2
					} else {
						escapeSeqLen = 1
					}
					groupName = groupName[0 : len(groupName)-escapeSeqLen]
					groupkey := group.GroupKeyShare{
						GroupName: groupName,
						Host:      gr.SelfId,
						Key:       "xxx",
					}
					group.CurrentGroupShareKey = &groupkey
					grouptopickey := groupkey.GenerateGroupKey()
					oldGroupRoom := gr
					oldGroupRoom.ExitRoom()
					time.Sleep(2 * time.Second)
					newGroupRoom, err := group.JoinGroup(gr.HostP2P, gr.UserName, grouptopickey)
					group.CurrentGroupRoom = newGroupRoom
					if err != nil {
						fmt.Println("Error while joining the new group")
						continue
					}
					gr = newGroupRoom
					fmt.Print("Waiting for queues to adapt")
					time.Sleep(1 * time.Second)
					fmt.Print(".")
					time.Sleep(1 * time.Second)
					fmt.Print(".")
					time.Sleep(1 * time.Second)
					fmt.Print(".")
					break
				case 2:
					fmt.Println("Enter the new user name")
					userName, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error during reading input from the SDI")
						continue
					}
					escapeSeqLen := 0
					if runtime.GOOS == "windows" {
						escapeSeqLen = 2
					} else {
						escapeSeqLen = 1
					}
					userName = userName[0 : len(userName)-escapeSeqLen]
					gr.UpdateUserName(userName)
					break
				case 3:
					fmt.Println("These are the list of peers in the Group " + gr.GroupName)
					for _, peer := range gr.PeerList() {
						fmt.Println(peer)
					}
					break
				case 4:
					fmt.Println("These are the list of your service peers currently active")
					for _, peer := range gr.HostP2P.Host.Network().Peers() {
						fmt.Println(peer)
					}
					break
				case 5:
					group.Peerlist = nil
					for id, peer := range gr.HostP2P.Host.Network().Peers() {
						group.Peerlist = append(group.Peerlist, group.ServicePeer{Id: id, PeerId: peer})
					}
					for _, servicepeer := range group.Peerlist {
						fmt.Println(servicepeer)
					}
					var choice int
					fmt.Scanln(&choice)
					var grphostid peer.ID
					for _, servicepeer := range group.Peerlist {
						if choice == servicepeer.Id {
							grphostid = servicepeer.PeerId
							break
						}
					}
					fmt.Println("Enter the group name to request to join:")
					groupName, err := reader.ReadString('\n')
					escapeSeqLen := 0
					if runtime.GOOS == "windows" {
						escapeSeqLen = 2
					} else {
						escapeSeqLen = 1
					}
					groupName = groupName[0 : len(groupName)-escapeSeqLen]
					if err != nil {
						fmt.Println("Error during reading the group name")
						continue
					}
					fmt.Println("Tryin to enter the group [" + groupName + "]")
					fmt.Println("Enter the message to the group host :")
					message, err := reader.ReadString('\n')
					message = message[0 : len(message)-escapeSeqLen]
					if err != nil {
						fmt.Println("Error during reading the message for group request")
						continue
					}
					group.GroupJoinRequest(ctx, host, groupName, grphostid, message)
					fmt.Println("Done with the group join request")
					break
				case 6:
					group.Peerlist = nil
					for id, peer := range gr.HostP2P.Host.Network().Peers() {
						group.Peerlist = append(group.Peerlist, group.ServicePeer{Id: id, PeerId: peer})
					}
					for _, servicepeer := range group.Peerlist {
						fmt.Println(servicepeer)
					}
					var choice int
					fmt.Scanln(&choice)
					var grphostid peer.ID
					for _, servicepeer := range group.Peerlist {
						if choice == servicepeer.Id {
							grphostid = servicepeer.PeerId
							break
						}
					}
					_, err = host.NewStream(ctx, grphostid, group.TestProtocol)
					if err != nil {
						fmt.Println("[ERROR] - in establishing a test stream")
					} else {
						fmt.Println("[SUCCESS] - in establishing a test stream")
					}
					break
				case 7:
					fmt.Println("************ACTIVEPEERTABLE*****************")
					group.PrintPeerTable()
					break
				default:
					fmt.Println("Bad command")

				}
			} else {
				escapeSeqLen := 0
				if runtime.GOOS == "windows" {
					escapeSeqLen = 2
				} else {
					escapeSeqLen = 1
				}
				msg := input[0 : len(input)-escapeSeqLen]
				chatmsg := &group.Chatmessage{
					Message:    msg,
					SenderID:   gr.SelfId,
					SenderName: gr.UserName,
				}
				chatbytes, err := json.Marshal(chatmsg)
				if err != nil {
					fmt.Println("[ERROR] - during marhsalling chat message")
					continue
				}
				msgpacket := &group.Packet{
					PacketType: "<chat>",
					Content:    chatbytes,
				}
				// fmt.Println("sending message to outbound queue")
				go func() {
					gr.Outbound <- *msgpacket
					// fmt.Println("Message sent to outbound queue")
				}()

			}
		}
	}
}
