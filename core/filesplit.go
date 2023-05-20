package core

import (
	"fmt"
	"os"
)

const BufferSize int = 8388608

var piecesMapList []PieceInfo

func (fI *FileInfo) SplitAndSave(file *os.File) {
	var folderName FolderName = FolderName(string(piecefolder) + "/" + fmt.Sprint(fI.FileName, "_", fI.FileType, "_"+fI.UniqueID.String()))
	folderName.MakeFolder()
	buffer := make([]byte, BufferSize)
	iterationCount := fI.FileSize / BufferSize
	lastPieceSize := BufferSize
	if fI.FileSize%BufferSize > 0 {
		iterationCount += 1
		lastPieceSize = fI.FileSize % BufferSize
	}
	lastPieceBuffer := make([]byte, lastPieceSize)
	fI.FilePieces = iterationCount
	fI.preparePieceDistributionList()
	for i := 0; i < iterationCount; i++ {
		if i == iterationCount-1 {
			file.Read(lastPieceBuffer)
			pieceFileName := fmt.Sprintf("part_%d.%s", i, fI.FileType)

			pieceFile, err := os.Create(string(folderName) + "/" + pieceFileName)
			if err != nil {
				fmt.Println("[ERROR] - during creating the piece file :", pieceFileName)
			} else {
				pieceFile.Write(lastPieceBuffer)
				pieceFile.Close()
				pieceInfo := ComposePieceInfo(pieceFileName, lastPieceSize, fI.FileName, fI.FileType, fI.UniqueID)
				pieceInfo.AddSources()
				fI.AppendPiecesMapList(pieceInfo)
				fmt.Println("Handled file ", pieceFileName)
				// fmt.Println(pieceInfo)
			}
		} else {
			file.Read(buffer)
			pieceFileName := fmt.Sprintf("part_%d.%s", i, fI.FileType)
			pieceFile, err := os.Create(string(folderName) + "/" + pieceFileName)
			if err != nil {
				fmt.Println("[ERROR] - during creating the piece file :", pieceFileName)
			} else {
				pieceFile.Write(buffer)
				pieceFile.Close()
				pieceInfo := ComposePieceInfo(pieceFileName, BufferSize, fI.FileName, fI.FileType, fI.UniqueID)
				pieceInfo.AddSources()
				fI.AppendPiecesMapList(pieceInfo)
				fmt.Println("Handled file ", pieceFileName)
				// fmt.Println(pieceInfo)
			}
		}

	}
	fI.Save()

}
