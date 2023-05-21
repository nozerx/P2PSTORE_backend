package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func HandleStreamPieceDownload(str network.Stream) {
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
	go pieceInfo.UploadPiece(str.Conn().RemotePeer())
}

func (p PieceInfo) UploadPiece(remotePeer peer.ID) {
	protocolName := "pieceUploadProtocol/" + p.ParentFileName + "/" + p.ParentFileType + "/" + p.ParentUniqueID.String() + "/" + p.PieceName + "/" + fmt.Sprint(p.PieceSize)
	str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(protocolName))
	if err != nil {
		fmt.Println("[ERROR] - during creating a new stream to protocol [" + protocolName + "]")
	} else {
		fmt.Println("[SUCCESS] - in establishing a stream to protocol [" + protocolName + "]")
		if err != nil {
			fmt.Println("[ERROR] - during retrieving piece size")
		} else {
			piecePath := string(uploadedPiecesFolder) + "/" + p.ParentFileName + "_" + p.ParentFileType + "_" + p.ParentUniqueID.String() + "/" + p.PieceName
			file, err := os.Open(piecePath)
			if err != nil {
				fmt.Println("[ERROR] - during trying to open piece [" + p.PieceName + "] at path [" + piecePath + "]")
			} else {
				buffer := make([]byte, p.PieceSize)
				file.Read(buffer)
				str.Write(buffer)
				fmt.Println("[SUCCESS] - sending piece to requester [" + remotePeer.String() + "]")
			}
		}

		str.Close()
	}
}
