package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	protocolName := "pieceDownloadProtocol/" + p.ParentFileName + "/" + p.ParentFileType + "/" + p.ParentUniqueID.String() + "/" + p.PieceName + "/" + fmt.Sprint(p.PieceSize)
	str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(protocolName))
	if err != nil {
		fmt.Println("[ERROR] - during creating a new stream to protocol [" + protocolName + "]")
	} else {
		fmt.Println("[SUCCESS] - in establishing a stream to protocol [" + protocolName + "]")
		folderName := FolderName(string(uploadedPiecesFolder) + "/" + fmt.Sprint(p.ParentFileName, "_", p.ParentFileType, "_"+p.ParentUniqueID.String()))
		err = folderName.MakeFolder()
		if err != nil {
			fmt.Println("[ERROR] - during creation of folder [" + folderName + "]")
			fmt.Println("[ERROR INFO] - ", err.Error())
			return
		} else {
			fmt.Println("Trying to recieve the piece [" + p.ParentFileName + "/" + p.PieceName + "]")
			file, err := os.Create(string(folderName) + "/" + p.PieceName)
			buffer := make([]byte, 1)
			reader := bufio.NewReader(str)
			if err != nil {
				fmt.Println("[ERROR] - during creation of recieving file [" + p.PieceName + "]")
			} else {
				for {
					_, err := reader.Read(buffer)
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
					file.Write(buffer)

				}
				fmt.Println("[SUCCESS] - piece [" + p.ParentFileName + "/" + p.PieceName + "] fully recieved")
				file.Close()
			}
		}
		str.Close()
	}
}
