package core

import (
	"context"

	"github.com/google/uuid"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type FolderName string
type FileName string

type PieceInfo struct {
	PieceName      string
	PieceType      string
	PieceSize      int
	Sources        []peer.ID
	Status         []bool
	ParentFileName string
	ParentFileType string
	ParentUniqueID uuid.UUID
}

type FileInfo struct {
	FileName      string
	FileType      string
	FileSize      int
	FilePieces    int
	UniqueID      uuid.UUID
	Pieces        []PieceInfo
	PeerTableCopy []peer.ID
}

type PickDistributionList struct {
	Peer  peer.ID
	Count int
}

type DynamicFilePieceHandleProtocol struct {
	protocolName string
	count        int
}

type FileBasicInfo struct {
	FileName string
	FileType string
	UniqueID uuid.UUID
}
type downloadStat struct {
	pieceName      string
	downloadStatus bool
}

type FileDownloadHanlde struct {
	File  FileInfo
	Stats []*downloadStat
}

type uploadStat struct {
	pieceName    string
	uploadStatus bool
}

type FileUploadHandle struct {
	File FileInfo
	Stat []*uploadStat
}

//constants

const rootFolder FolderName = "core"
const mapfilefolder FolderName = "core/mapfiles"
const piecefolder FolderName = "core/piecefolders"
const uploadedPiecesFolder FolderName = "core/uploaded"
const downloadedPiecesFolder FolderName = "core/downloaded"
const sendfolder FolderName = "core/send"
const recievefolder FolderName = "core/recieve"

const FilePieceUploadProtocol string = "rex/file/upload/piece"
const FilePieceDownloadProtocol string = "rex/file/download/piece"
const MapFileShareProtocol string = "rex/file/meta/share"

//variables

var NodeHostCtx struct {
	Host host.Host
	Ctx  context.Context
}

var ManageProtocolList []DynamicFilePieceHandleProtocol

var DownloadQueue []*FileDownloadHanlde
var UploadQueue []*FileUploadHandle
