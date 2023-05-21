package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (fl FileInfo) Send() {
	for _, piece := range fl.Pieces {
		go piece.SendPiece(fl)
		time.Sleep(time.Second * 1)
	}
}

func (p PieceInfo) NotifyAllSources(pieceBytes []byte) {
	for _, remotePeer := range p.Sources {
		str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(FilePieceUploadProtocol))
		if err != nil {
			fmt.Println("[ERROR] - during creating the stream [" + FilePieceUploadProtocol + "] to Peer[" + remotePeer.String() + "]")
		} else {
			fmt.Println("[SUCCESS] - during creating the stream [" + FilePieceUploadProtocol + "] to Peer[" + remotePeer.String() + "]")
			str.Write(pieceBytes)
			str.Close()
		}
		time.Sleep(time.Second * 1)
	}
}

func (p PieceInfo) SendPiece(fl FileInfo) {
	pieceBytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling piece [", p.PieceName, "]")
	}
	go p.NotifyAllSources(pieceBytes)
	pieceProtocol := "pieceDownloadProtocol/" + p.ParentFileName + "/" + p.ParentFileType + "/" + p.ParentUniqueID.String() + "/" + p.PieceName + "/" + fmt.Sprint(p.PieceSize)
	// pieceProtocol := "pieceDownloadProtocol/" + p.PieceName + "/" + p.PieceType
	NodeHostCtx.Host.SetStreamHandler(protocol.ID(pieceProtocol), HandleStreamPieceDataUpload)
	ManageProtocolList = append(ManageProtocolList, DynamicFilePieceHandleProtocol{
		protocolName: pieceProtocol,
		count:        3,
	})
	fmt.Println("protocol [" + pieceProtocol + "] is commissioned for use")

}

func ManageStreamHandlerRemoval(protocolName string) {
	for i, proto := range ManageProtocolList {
		if strings.EqualFold(protocolName, proto.protocolName) {
			ManageProtocolList[i].count--
			fmt.Println("[INFO] - count decremented for protocol ["+protocolName+"]", ManageProtocolList[i].count)
			if ManageProtocolList[i].count == 0 {
				NodeHostCtx.Host.RemoveStreamHandler(protocol.ID(protocolName))
				fmt.Println("protocol [" + protocolName + "] is dissmissed after use")
			}
			break
		}
	}
}

func HandleStreamPieceDataUpload(str network.Stream) {
	protocolName := string(str.Protocol())
	fileName := strings.Split(protocolName, "/")[1]
	fileType := strings.Split(protocolName, "/")[2]
	uniqueId := strings.Split(protocolName, "/")[3]
	pieceName := strings.Split(protocolName, "/")[4]
	pieceSize, err := strconv.Atoi(strings.Split(protocolName, "/")[5])
	if err != nil {
		fmt.Println("[ERROR] - during retrieving piece size")
		fmt.Println(strings.Split(protocolName, "/")[5])
		return
	}
	fmt.Println("An incomming strem requesting file [" + fileName + "/" + pieceName + "] identified")
	piecePath := string(piecefolder) + "/" + fmt.Sprint(fileName, "_", fileType, "_"+uniqueId) + "/" + pieceName
	SendPiece(str, piecePath, pieceSize)
	str.Close()
	ManageStreamHandlerRemoval(protocolName)
	// NodeHostCtx.Host.RemoveStreamHandler(protocol.ID(protocolName))
	// fmt.Println("protocol [" + protocolName + "] is dissmissed after use")
}

func SendPiece(str network.Stream, path string, size int) {
	fmt.Println("Piece Sending Started..")
	file, err := os.Open(path)
	reader := bufio.NewReader(file)
	buffer := make([]byte, size)
	if err != nil {
		fmt.Println("[ERROR] - during opening the file for sending")
		return
	} else {
		reader.Read(buffer)
		_, err = str.Write(buffer)
		if err != nil {
			fmt.Println("[ERROR] - during sending the piece to RemotePeer [" + str.Conn().RemotePeer().String() + "]")
			return
		} else {
			fmt.Println("[SUCCESS] - in sending the piece [" + path + "] to peer [" + str.Conn().RemotePeer().String() + "]")
			str.Close()
		}

	}

}
