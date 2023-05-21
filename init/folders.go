package init

import (
	"fmt"
	"os"
)

// Here we define type for foldername so we can define methods on this type
// So we can call the methods (MakeFolder()) on any of the variables of this type
type folderName string

// Here we define the type for filenames so we can define methods on this type
// So we can call the methods () on any of the variable of this type
type FileName string

// Here we define all the folder needed for the node to initalize properly
const rootFolder folderName = "core"
const Nodeinfofolder folderName = "core/nodeinfo"
const Mapfilefolder folderName = "core/mapfiles"
const Piecefolder folderName = "core/piecefolders"
const UploadedPiecesFolder folderName = "core/uploaded"
const DownloadedPiecesFolder folderName = "core/downloaded"
const Sendfolder folderName = "core/send"
const Recievefolder folderName = "core/recieve"

// this list makes it easier to cycle through the file names and handle any errors encountered while creating them
var initfoldernames []folderName = []folderName{rootFolder, Nodeinfofolder, Mapfilefolder, Piecefolder, UploadedPiecesFolder, DownloadedPiecesFolder, Sendfolder, Recievefolder}

// This function handles creation of any new folders
func (fn folderName) MakeFolder() error {
	_, err := os.Stat(string(fn)) // checks if the folder already exists
	if os.IsNotExist(err) {       // if the folder doesn't exist
		fmt.Println("Creating the folder " + string(fn))
		err := os.Mkdir(string(fn), 0755) // creating the folder
		if err != nil {
			fmt.Println("[ERROR] during creating the folder" + string(fn))
			return err
		}
		fmt.Println("[SUCCESS] - in creating the folder [" + string(fn) + "]")
		return nil
	} else { // if the folder already exists
		fmt.Println("Folder " + string(fn) + " already exits")
		return nil
	}
}

// This function initializes all the core folders needed for the node to work as intended
func InitFolders() error {
	for _, folder := range initfoldernames {
		err := folder.MakeFolder()
		if err != nil {
			fmt.Println("[CAUTION] - node maynot function properly without all folders initialized")
			fmt.Println("[SUGGESSION] - Try restarting the node to solve this problem")
			return fmt.Errorf("[ERROR] - connot function properly without core folders initialized")
		}
	}
	return nil
}
