package filehandling

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
)

func HandleStreamFileShareMetaIncomming(str network.Stream) {
	buff = nil
	fmt.Println("File Recieve identified")
	streamReader := bufio.NewReader(str)
	buffer := make([]byte, 1)

	for {
		_, err := streamReader.Read(buffer)
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
		buff = append(buff, buffer...)

	}
	fmt.Println("Exited the loop")
	file := &File{}
	err := json.Unmarshal(buff, file)
	if err != nil {
		fmt.Println("[ERROR] - during unmarshalling")
	} else {
		str.Close()
		fmt.Println(file)
		go file.Recieve()
	}
}
