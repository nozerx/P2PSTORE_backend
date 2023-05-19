package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (fl FileInfo) Send() {
	for _, piece := range fl.Pieces {
		go piece.SendPiece(fl)
	}
}

func (p PieceInfo) GetRandomSource() peer.ID {
	return p.Sources[0]
}

func (p PieceInfo) SendPiece(fl FileInfo) {
	pieceBytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling piece [", p.PieceName, "]")
	}
	str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, p.GetRandomSource(), protocol.ID(FilePieceUploadProtocol))
	str.Write(pieceBytes)
	str.Close()
	// pieceProtocol := "pieceDownloadProtocol/" + fl.FileName + "/" + fl.UniqueID.String() + "/" + p.PieceName
	pieceProtocol := "pieceDownloadProtocol/" + p.PieceName + "/" + p.PieceType
	NodeHostCtx.Host.SetStreamHandler(protocol.ID(pieceProtocol), HandleStreamPieceDataUpload)
	fmt.Println("protocol [" + pieceProtocol + "] is commissioned for use")

}

func HandleStreamPieceDataUpload(str network.Stream) {
	protocolName := string(str.Protocol())
	pieceName := strings.Split(protocolName, "/")[1]
	fmt.Println("An incomming strem requesting file [" + pieceName + "] identified")
	str.Close()
	NodeHostCtx.Host.RemoveStreamHandler(protocol.ID(protocolName))
	fmt.Println("protocol [" + protocolName + "] is dissmissed after use")
}
