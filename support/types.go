package support

import (
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PieceInfo struct {
	PieceName      string
	PieceType      string
	PieceSize      int
	Sources        []peer.ID
	ParentFileName string
	ParentFileType string
	ParentUniqueID uuid.UUID
}

type FileInfo struct {
	FileName   string
	FileType   string
	FileSize   int
	FilePieces int
	UniqueID   uuid.UUID
	Pieces     []PieceInfo
}
