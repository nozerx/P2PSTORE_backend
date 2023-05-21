package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"p2pstore/support"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func parseUUID(uuidStr string) uuid.UUID {
	uuidObj, err := uuid.Parse(uuidStr)
	if err != nil {
		fmt.Println("[ERROR] - during parsing a uuid from the string")
		fmt.Println("[" + uuidStr + "]")
		return uuid.Nil
	}
	return uuidObj
}

func ListAvailableFiles() []FileBasicInfo {
	files, err := ioutil.ReadDir(string(mapfilefolder))
	if err != nil {
		fmt.Println("[ERROR] - during reading files available in the node")
		return nil
	}
	var availableFileList []FileBasicInfo = nil
	for _, file := range files {
		filebasicinfo := FileBasicInfo{
			FileName: strings.Split(file.Name(), "_")[0],
			FileType: strings.Split(file.Name(), "_")[1],
			UniqueID: parseUUID(strings.Split(strings.Split(file.Name(), "_")[2], ".")[0]),
		}
		availableFileList = append(availableFileList, filebasicinfo)

	}
	return availableFileList

}

func GetMapFileSize(fileNamePath string) int {
	file, err := os.Stat(fileNamePath)
	if err != nil {
		fmt.Println("[ERROR] - during determining the file size")
		return 0
	} else {
		return int(file.Size())
	}
}

func (fl FileBasicInfo) HandleFileDownload(fileUniqueID string) error {
	filepath := string(mapfilefolder) + "/" + fl.FileName + "_" + fl.FileType + "_" + fl.UniqueID.String() + ".txt"
	bufferSize := GetMapFileSize(filepath)
	buffer := make([]byte, bufferSize)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("[ERROR] - during trying to open the file map list")
		fmt.Println("[PATH : " + filepath + "]")
		return nil
	}
	file.Read(buffer)
	fmt.Println()

	// err = json.Unmarshal(buffer, &fileInfo)
	// if err != nil {
	// 	fmt.Println("[ERROR] - during unmarhsalling the file map list")
	// 	fmt.Println(err)
	// 	// return nil
	// }
	fileInfoSupport := support.GetFileInfo(buffer)
	var fileInfo FileInfo = FileInfo{
		FileName:   fileInfoSupport.FileName,
		FileType:   fileInfoSupport.FileType,
		FileSize:   fileInfoSupport.FileSize,
		FilePieces: fileInfoSupport.FilePieces,
		UniqueID:   fileInfoSupport.UniqueID,
		Pieces:     nil,
	}
	for _, piece := range fileInfoSupport.Pieces {
		fileInfo.Pieces = append(fileInfo.Pieces, PieceInfo{
			PieceName:      piece.PieceName,
			PieceType:      piece.PieceType,
			PieceSize:      piece.PieceSize,
			Sources:        piece.Sources,
			ParentFileName: piece.ParentFileName,
			ParentFileType: piece.ParentFileType,
			ParentUniqueID: piece.ParentUniqueID,
		})
	}
	fmt.Println("====================================================")
	fmt.Println("FILE NAME : [", fileInfo.FileName, "]")
	fmt.Println("FILE TYPE : [", fileInfo.FileType, "]")
	fmt.Println("FILE SIZE: [", fileInfo.FileSize, "]")
	fmt.Println("FILE PIECES COUNT : [", fileInfo.FilePieces, "]")
	fmt.Println("FILE UNIQUE-ID : [", fileInfo.UniqueID, "]")
	fmt.Println("----------------------------------------------------")
	for _, piece := range fileInfo.Pieces {
		fmt.Println("[PIECE] :["+piece.PieceName+"."+piece.PieceType+",", piece.Sources, "]")
	}
	fmt.Println("----------------------------------------------------")
	fmt.Println("====================================================")
	fileInfo.Retrieve()
	return nil
}

func (fl FileInfo) Retrieve() {
	var fileDownloadHanlder FileDownloadHanlde = FileDownloadHanlde{
		File:  fl,
		Stats: nil,
	}
	DownloadQueue = EnqueuDownload(&fileDownloadHanlder)
	for _, piece := range fl.Pieces {
		go piece.GetPiece()
		fileDownloadHanlder.Stats = append(fileDownloadHanlder.Stats, &downloadStat{
			pieceName:      piece.PieceName,
			downloadStatus: false,
		})
		time.Sleep(time.Second * 1)
	}
	CheckDownloadStatus(fl)
}

func (p PieceInfo) TryAllSources() network.Stream {
	for _, remotePeer := range p.Sources {
		str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(FilePieceDownloadProtocol))
		if err != nil {
			fmt.Println("[INFO] - Source [" + remotePeer.String() + "] not active")
			continue
		} else {
			fmt.Println("[INFO] - getting [" + p.PieceName + "] from source [" + remotePeer.String() + "]")
			return str
		}
	}
	fmt.Println("[ERROR] - file piece [" + p.PieceName + "] cannot be sourced")
	return nil
}

func (p PieceInfo) GetPiece() {
	str := p.TryAllSources()
	if str != nil {
		err := p.SendPieceMetaData(str)
		if err == nil {
			dynamicProtocol := "pieceUploadProtocol/" + p.ParentFileName + "/" + p.ParentFileType + "/" + p.ParentUniqueID.String() + "/" + p.PieceName + "/" + fmt.Sprint(p.PieceSize)
			NodeHostCtx.Host.SetStreamHandler(protocol.ID(dynamicProtocol), HandleStreamDownloadPiece)
		}
	} else {
		fmt.Println("[ABORT] - file download aborted piece missing | [" + p.PieceName + "]")
		return
	}
}

func (p PieceInfo) SendPieceMetaData(str network.Stream) error {
	buff, err := json.Marshal(p)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling")
		return err
	} else {
		str.Write(buff)
		str.Close()
		return nil
	}
}

func HandleStreamDownloadPiece(str network.Stream) {
	// fmt.Println("Piece Recieved", str.Protocol())
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
	} else {
		fmt.Println("An incomming strem sending file [" + fileName + "/" + pieceName + "] identified")
		folderName := FolderName(string(downloadedPiecesFolder) + "/" + fmt.Sprint(fileName, "_", fileType, "_"+uniqueId))
		err = folderName.MakeFolder()
		if err != nil {
			fmt.Println("[ERROR] - during creation of folder [" + folderName + "]")
			fmt.Println("[ERROR INFO] - ", err.Error())
			return
		} else {
			fmt.Println("Trying to recieve the piece [" + fileName + "/" + pieceName + "]")
			file, err := os.Create(string(folderName) + "/" + pieceName)
			buffer := make([]byte, 1)
			reader := bufio.NewReader(str)
			if err != nil {
				fmt.Println("[ERROR] - during creation of recieving file [" + pieceName + "]")
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
				fmt.Println("[SUCCESS] - piece [" + fileName + "/" + pieceName + "] fully recieved")
				// creating a dummy piece info object to set downloaded status
				piece := PieceInfo{
					PieceName:      pieceName,
					PieceType:      ":",
					PieceSize:      pieceSize,
					Sources:        []peer.ID{},
					ParentFileName: fileName,
					ParentFileType: fileType,
					ParentUniqueID: uuid.Nil,
				}

				piece.SetDownloadedStatus()
				file.Close()
			}
		}
		str.Close()
	}

}

func (p PieceInfo) SetDownloadedStatus() {
	for _, fileHandle := range DownloadQueue {
		if fileHandle.File.FileName == p.ParentFileName {
			for _, piece := range fileHandle.Stats {
				if piece.pieceName == p.PieceName {
					fmt.Println("[INFO] - piece [" + p.ParentFileName + "/" + p.PieceName + "] is set to downloaded")
					piece.downloadStatus = true
				}
			}
		}
	}
}

func CheckDownloadStatus(file FileInfo) {
	var fh *FileDownloadHanlde
	for _, fileHandle := range DownloadQueue {
		if fileHandle.File.FileName == file.FileName {
			fh = fileHandle
		}
	}
	var fullyDownloaded bool
	for {
		time.Sleep(time.Second * 1)
		fullyDownloaded = true
		for _, piece := range fh.Stats {
			if piece.downloadStatus == false {
				fullyDownloaded = false
			}
		}
		if fullyDownloaded == true {
			fmt.Println("[INFO] - all pieces of file [" + file.FileName + "." + file.FileType + "] downloaded")
			fmt.Println("[INFO] - tryig to merge the pices to get the intended file")
			file.Merge()
			return
		}
	}
}

func (p *PieceInfo) GetPieceBytes(fI *FileInfo) []byte {
	piecePath := string(downloadedPiecesFolder) + "/" + fI.FileName + "_" + fI.FileType + "_" + fI.UniqueID.String() + "/" + p.PieceName
	file, err := os.Open(piecePath)
	if err != nil {
		fmt.Println("[ERROR] - during opening the piece ", p.PieceName)
		fmt.Println("[PATH:" + piecePath + "]")
		return nil
	}
	pieceBuffer := make([]byte, p.PieceSize)
	file.Read(pieceBuffer)
	return pieceBuffer
}

func (fI *FileInfo) Merge() {
	destinationFilePath := string(recievefolder) + "/" + string(fI.FileName) + "." + string(fI.FileType)
	file, err := os.Create(destinationFilePath)
	if err != nil {
		fmt.Println("[ERROR] - during creating the destination file")
		fmt.Println("[PATH :" + destinationFilePath + "]")
		return
	}
	for i := 0; i < fI.FilePieces; i++ {
		file.Write(fI.Pieces[i].GetPieceBytes(fI))
	}
	fmt.Println("[SUCCESS] - in reassembling the file")
	file.Close()
}

func EnqueuDownload(fileHandle *FileDownloadHanlde) []*FileDownloadHanlde {
	DownloadQueue = append(DownloadQueue, fileHandle)
	return DownloadQueue
}
