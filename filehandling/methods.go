package filehandling

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func HandleStreamFileShare(str network.Stream) {
	file, FileSendQueue := FileSendQueue.Dequeue()
	fmt.Println(FileSendQueue)
	fmt.Println(file)
	go file.SendFile(str)
}

func CreatNewFileObject(fileName string, fileType string, host host.Host) (*File, error) {
	file := &File{
		FileName:  fileName,
		FileType:  fileType,
		FileSize:  0,
		FileOwner: host.ID(),
	}
	size, err := file.GetFileSize()
	if err != nil {
		fmt.Println("[ERROR] - during determining the file size")
		return nil, err
	}
	file.FileSize = size
	return file, nil
}

func (fl File) GetFileSize() (int, error) {
	file, err := os.Stat(string(sendfolder) + "/" + fl.FileName)
	if err != nil {
		fmt.Println("[ERROR] - during determining the file size")
		return 0, err
	} else {
		return int(file.Size()), nil
	}
}

func (fl File) SendFile(str network.Stream) {
	fmt.Println("File Sending Started..")
	file, err := os.Open(string(sendfolder) + "/" + fl.FileName)
	if err != nil {
		fmt.Println("[ERROR] - during opening the file for sending")
		return
	} else {
		if err != nil {
			fmt.Println("[ERROR] - during creating a new stream on protocol [", FileShareProtocol, "]")
			return
		}
		fl.toStream(str, file)
		// err := file.Close()
		// if err != nil {
		// 	fmt.Println("[ERROR] - during closing the file")
		// }
		// str.Close()
	}
}

func (fl File) toStream(str network.Stream, file *os.File) {
	fmt.Println("Inside to Stream")
	iterationcount := fl.FileSize / buffSize
	buffer := make([]byte, buffSize)
	streamWriter := bufio.NewWriter(str)
	fmt.Println(iterationcount)
	for i := 0; i < iterationcount; i++ {
		fmt.Println("iteration")
		_, err := file.Read(buffer)
		if err == io.EOF {
			fmt.Println("File send to stream completely")
			break
		}
		if err != nil {
			fmt.Println("Error while reading from the file")
		}
		sendByte, err := streamWriter.Write(buffer)
		if err != nil {
			fmt.Println("Error while sending the buffer to the stream")
		} else {
			err := streamWriter.Flush()
			if err != nil {
				fmt.Println("Error while flushing")
			}
			fmt.Println("Send ", sendByte, " bytes to stream")
		}

	}
	fmt.Println("Outside the loop")
	leftByte := fl.FileSize % buffSize
	additionalBuffer := make([]byte, leftByte)
	_, err := file.Read(additionalBuffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("File send to stream completely")
		} else {
			fmt.Println("Error while reading the last buffer")

		}
	}
	sendByte, err := streamWriter.Write(additionalBuffer)
	if err != nil {
		fmt.Println("Error while sending the buffer to the stream")
	} else {
		err := streamWriter.Flush()
		if err != nil {
			fmt.Println("Error while flushing")
		}
		fmt.Println("Send ", sendByte, " bytes to stream")
	}
	fmt.Println("End of toStream function")
}

func (fl File) Recieve() {
	file, err := os.Create(string(recievefolder) + "/" + fl.FileName)
	if err != nil {
		fmt.Println("[ERROR] - during trying to create the recieving file")
	} else {
		str, err := NodeHostCtx.Host.NewStream(NodeHostCtx.Ctx, fl.FileOwner, protocol.ID(FileShareProtocol))
		if err != nil {
			fmt.Println("[ERROR] during creating a new stream on protocol [", FileShareProtocol, "]")
		} else {
			fl.fromStream(str, file)
			err := file.Close()
			if err != nil {
				fmt.Println("[ERROR] - during closing the file")
			}
			str.Close()
		}
	}
}

func (fl File) fromStream(str network.Stream, file *os.File) {

	streamReader := bufio.NewReader(str)
	buffer := make([]byte, bufferSize)
	iterationCount := fl.FileSize / bufferSize
	fmt.Println(iterationCount)
	for i := 0; i < iterationCount; i++ {
		_, err := streamReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file recieved")
			} else {
				fmt.Println("Error while reading file bytes from stream")
			}
		} else {
			_, err := file.Write(buffer)
			if err != nil {
				fmt.Println("Error while writing to the recieving file")
			}
		}
	}
	leftByte := fl.FileSize % bufferSize
	additonalBuffer := make([]byte, leftByte)
	_, err := streamReader.Read(additonalBuffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("End of file recieved")
		} else {
			fmt.Println("Error while reading file bytes from stream for last piece")
		}
	} else {
		_, err := file.Write(additonalBuffer)
		if err != nil {
			fmt.Println("Error while writing to the recieving file")
		} else {
			fmt.Println("Completely recieved the file")
		}
	}
}
