package core

import (
	"fmt"
	"io"
	"os"
	"p2pstore/support"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (fl FileInfo) SendMapFile(remotePeer peer.ID) {
	path := string(mapfilefolder) + "/" + fl.FileName + "_" + fl.FileType + "_" + fl.UniqueID.String() + ".txt"
	bufferSize := GetMapFileSize(path)
	buffer := make([]byte, bufferSize)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[ERROR][MAPFILE-SHARE] - error during opening the mapfile for sharing")
	} else {
		_, err := file.Read(buffer)
		if err != nil {
			fmt.Println("[ERROR][MAPFILE-SHARE] - during reading from map file")
		} else {
			if remotePeer == peer.ID(fmt.Sprint("nil")) {
				for _, peer := range fl.PeerTableCopy {
					go sendTo(file, peer, buffer)
				}
			} else {
				sendTo(file, remotePeer, buffer)
			}
		}

	}
	file.Close()
}

func sendTo(file *os.File, remotePeer peer.ID, buffer []byte) {
	str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, remotePeer, protocol.ID(MapFileShareProtocol))
	if err != nil {
		fmt.Println("[ERROR][MAPFILE-SHARE] - during establishing a stream to peer [" + remotePeer.Pretty() + "]")
	} else {

		_, err := str.Write(buffer)
		if err != nil {
			fmt.Println("[ERROR][MAPFILE-SHARE] - during writing to stream")
		} else {
			str.Close()
		}
	}
}

func HandleStreamMapFileShare(str network.Stream) {
	var buff []byte = nil
	buffer := make([]byte, 1)
	for {
		_, err := str.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Fully Recieved")
				break
			}
			fmt.Println("[ERROR][MAPFILE-SHARE] - during reading from the input stream")
			fmt.Println(err.Error())
			fmt.Println("[ABORT][MAPFILE-SHARE] - handling mapfile share from " + str.Conn().RemotePeer())
			return
		}
		buff = append(buff, buffer...)
	}
	mapfile := support.GetFileInfo(buff)
	path := string(mapfilefolder) + "/" + mapfile.FileName + "_" + mapfile.FileType + "_" + mapfile.UniqueID.String() + ".txt"
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("[ERROR][FILE-CREATION][MAPFILE-SHARE] - during creation of recieving file for mapfile")
	} else {
		_, err := file.Write(buff)
		if err != nil {
			fmt.Println("[ERROR][WRITING][MAPFILE-SHARE] - during writing to recieving file for mapfile")
		} else {
			fmt.Println("[SUCCESS][MAPFILE-SHARE] - during recieving mapfile")
		}
	}
	file.Close()

}
