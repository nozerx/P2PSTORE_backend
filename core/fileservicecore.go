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
		input, err := reader.ReadString('\n')
		escapeSeqLen := 0
		if runtime.GOOS == "windows" {
			escapeSeqLen = 2
		} else {
			escapeSeqLen = 1
		}
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		if err != nil {
			fmt.Println("[ERROR] - during reading file name")
			break
		}
		var fileName FileName = FileName(input[:len(input)-escapeSeqLen])
		fmt.Println("The Entered file name is :" + fileName)
		fileName.HandleFile()
		break
	case 2:
		break
	case 3:
		break
	default:

	}

}
