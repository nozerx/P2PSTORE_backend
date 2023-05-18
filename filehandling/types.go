package filehandling

import (
	"context"

	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type File struct {
	FileName  string
	FileType  string
	FileSize  int
	FileOwner peer.ID
	UniqueID  uuid.UUID
}
type folderName string
type FileName string
type Queue []File

//constants

const rootFolder folderName = "core"
const mapfilefolder folderName = "core/mapfiles"
const piecefolder folderName = "core/piecefolders"
const sendfolder folderName = "core/send"
const recievefolder folderName = "core/recieve"
const FileShareProtocol string = "rex/file/share"
const FileShareMetaDataProtocol string = "rex/file/share/metadata"
const buffSize = 10485760
const bufferSize = 1

//variables

var buff []byte
var NodeHostCtx struct {
	Host host.Host
	Ctx  context.Context
}

var FileSendQueue Queue
