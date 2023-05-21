package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func fileServiceMenu(reader *bufio.Reader) {
	fmt.Println("********************[FILE MENU]*******************")
	fmt.Println("[1]-File Upload\n[2]-File Download\n[3]-File Meta Share")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		fmt.Println("---------------------[FILE UPLOAD]--------------------")
		fmt.Println("Enter the File Name to Be Uploaded :")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] - during reading file name")
			break
		}
		escapeSeqLen := 0
		if runtime.GOOS == "windows" {
			escapeSeqLen = 2
		} else {
			escapeSeqLen = 1
		}
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

		var fileName FileName = FileName(input[:len(input)-escapeSeqLen])
		fmt.Println("The Entered file name is :" + fileName)
		fileName.HandleFile()
		break
	case 2:
		fmt.Println("---------------------[FILE DOWNLOAD]--------------------")
		fmt.Println("-------AVAILABLE FILES---------")
		for i, files := range ListAvailableFiles() {
			fmt.Println("[", i, "]-["+files.FileName+"."+files.FileType+"]")
		}
		fmt.Println("-------------------------------")
		fmt.Print(">")
		var chooseFile int
		fmt.Scanln(&chooseFile)
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		file := ListAvailableFiles()[chooseFile]
		var fileName FileName = FileName(file.FileName + "." + file.FileType)
		fmt.Println("The Entered file name is :" + fileName)
		err := file.HandleFileDownload(file.UniqueID.String())
		if err != nil {
			fmt.Println("[ABORTED] - file download")
		}
		break
	case 3:
		break
	default:

	}

}
