package support

import (
	"encoding/json"
	"fmt"
)

func GetFileInfo(buff []byte) FileInfo {
	file := FileInfo{}
	err := json.Unmarshal(buff, &file)
	if err != nil {
		fmt.Println(err)
		return FileInfo{}
	}
	return file
}
