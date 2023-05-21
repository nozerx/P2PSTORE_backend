package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"p2pstore/support"
	"strings"

	"github.com/google/uuid"
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
	return nil
}
