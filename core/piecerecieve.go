package core

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func HandleStreamPieceUpload(str network.Stream) {
	var buff []byte = nil
	fmt.Println("File Recieve identified")
	buffer := make([]byte, 1)

	for {
		_, err := str.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Fully Recieved")
				break
			}
			fmt.Println("[ERROR] - during reading from the input stream")
			fmt.Println(err.Error())
			fmt.Println("[ABORT] - handling join request from " + str.Conn().RemotePeer())
			return
		}
		buff = append(buff, buffer...)
	}
	pieceInfo := &PieceInfo{}
	json.Unmarshal(buff, pieceInfo)
	fmt.Println(pieceInfo.PieceName, pieceInfo.PieceSize)
	go pieceInfo.RecievePiece(str.Conn().RemotePeer())
}

func (p PieceInfo) RecievePiece(remotePeer peer.ID) {
	protocolName := "pieceDownloadProtocol/" + p.PieceName + "/" + p.PieceType
	str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(protocolName))
	if err != nil {
		fmt.Println("[ERROR] - during creating a new stream to protocol [" + protocolName + "]")
	} else {
		fmt.Println("[SUCCESS] - in establishing a stream to protocol [" + protocolName + "]")
		str.Close()
	}
}
