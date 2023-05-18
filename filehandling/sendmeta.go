package filehandling

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (fq Queue) Enqueue(fl File) Queue {
	fq = append(fq, fl)
	return fq
}

func (fq Queue) Dequeue() (File, Queue) {
	file := fq[0]
	fq = fq[1:]
	return file, fq
}

func (fl File) SendMeta(ctx context.Context, host host.Host, remotePeer peer.ID) error {

	metaBytes, err := json.Marshal(fl)
	if err != nil {
		fmt.Println("[ERROR] - during marshalling")
		return err
	} else {
		str, err := host.NewStream(ctx, remotePeer, protocol.ID(FileShareMetaDataProtocol))
		if err != nil {
			fmt.Println("[ERROR] -during creating a new stream on protocol [", FileShareMetaDataProtocol, "]")
		} else {
			FileSendQueue = FileSendQueue.Enqueue(fl)
			str.Write(metaBytes)
			str.Close()
		}
	}
	return nil
}
