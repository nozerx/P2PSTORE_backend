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
	var fileUploadHandle FileUploadHandle = FileUploadHandle{
		File: fl,
		Stat: nil,
	}
	UploadQueue = EnqueueUpload(&fileUploadHandle)
	for _, piece := range fl.Pieces {
		go piece.SendPiece(fl)
		fileUploadHandle.Stat = append(fileUploadHandle.Stat, &uploadStat{
			pieceName:    piece.PieceName,
			uploadStatus: false,
		})
		time.Sleep(time.Second * 1)
	}
	go CheckUploadStatus(&fl)
}

func (p *PieceInfo) NotifyAllSources(pieceBytes []byte) {
	for i, remotePeer := range p.Sources {
		str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(FilePieceUploadProtocol))
		if err != nil {
			fmt.Println("[ERROR] - during creating the stream [" + FilePieceUploadProtocol + "] to Peer[" + remotePeer.String() + "]")
		} else {
			fmt.Println("[SUCCESS] - during creating the stream [" + FilePieceUploadProtocol + "] to Peer[" + remotePeer.String() + "]")
			p.Status[i] = true
			str.Write(pieceBytes)
			str.Close()
		}
		time.Sleep(time.Second * 1)
	}
}

func (p *PieceInfo) SendPiece(fl FileInfo) {
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
	fmt.Println("[INFO] - protocol [" + pieceProtocol + "] is commissioned for use")

}

func ManageStreamHandlerRemoval(protocolName string) {
	for i, proto := range ManageProtocolList {
		if strings.EqualFold(protocolName, proto.protocolName) {
			ManageProtocolList[i].count--
			// fmt.Println("[INFO] - count decremented for protocol ["+protocolName+"]", ManageProtocolList[i].count)
			if ManageProtocolList[i].count == 0 {
				NodeHostCtx.Host.RemoveStreamHandler(protocol.ID(protocolName))
				fmt.Println("[INFO] - protocol [" + protocolName + "] is dissmissed after use")
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
	// fmt.Println("An incomming strem requesting file [" + fileName + "/" + pieceName + "] identified")
	piecePath := string(piecefolder) + "/" + fmt.Sprint(fileName, "_", fileType, "_"+uniqueId) + "/" + pieceName
	SendPiece(str, piecePath, pieceSize, fileName, pieceName)
	str.Close()
	ManageStreamHandlerRemoval(protocolName)
	// NodeHostCtx.Host.RemoveStreamHandler(protocol.ID(protocolName))
	// fmt.Println("protocol [" + protocolName + "] is dissmissed after use")
}

func SendPiece(str network.Stream, path string, size int, fileName string, pieceName string) {
	// fmt.Println("Piece Sending Started..")
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
			SetPieceUploadStatus(fileName, pieceName)
			str.Close()
		}

	}

}

func SetPieceUploadStatus(FileName string, PieceName string) {
	for _, fileHandle := range UploadQueue {
		if fileHandle.File.FileName == FileName {
			for _, piece := range fileHandle.Stat {
				if piece.pieceName == PieceName {
					// fmt.Println("[INFO] - piece [" + FileName + "/" + PieceName + "] is set to uploaded")
					piece.uploadStatus = true
				}
			}
		}
	}
}

func UpdateUploadStatus(file FileInfo) {
	path := string(mapfilefolder) + "/" + file.FileName + "_" + file.FileType + "_" + file.UniqueID.String() + ".txt"
	fl, err := os.Open(path)
	if err != nil {
		fmt.Println("[ERROR] - during opening the map file to update it")
	} else {
		buff, err := json.Marshal(file)
		if err != nil {
			fmt.Println("[ERROR] - during marshalling mapfile during update")
		} else {
			fl.Write(buff)
			fl.Close()
		}
	}

}

func EnqueueUpload(fileHandle *FileUploadHandle) []*FileUploadHandle {
	UploadQueue = append(UploadQueue, fileHandle)
	return UploadQueue
}

func CheckUploadStatus(file *FileInfo) {
	fmt.Println("[INFO] - fileUpload Handler for [" + file.FileName + "." + file.FileType + "] started")
	var fh *FileUploadHandle
	for _, fileHandle := range UploadQueue {
		if fileHandle.File.FileName == file.FileName {
			fh = fileHandle
		}
	}
	var fullyUploaded bool
	for {
		time.Sleep(time.Second * 1)
		fullyUploaded = true
		for _, piece := range fh.Stat {
			if piece.uploadStatus == false {
				fullyUploaded = false
			}
		}
		if fullyUploaded == true {
			fmt.Println("[INFO] - all pieces of file [" + file.FileName + "." + file.FileType + "] uploaded")
			fmt.Println("[INFO] - fileUpload Handler for [" + file.FileName + "." + file.FileType + "] ended")
			UpdateUploadStatus(*file)
			return
		}
	}
}
