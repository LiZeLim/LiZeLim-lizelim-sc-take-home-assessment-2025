package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	res := folder.GetAllFolders()
	

	// example usage
	folderDriver := folder.NewDriver(res)
	orgFolder := folderDriver.GetFoldersByOrgID(orgID)

	folder.PrettyPrint(res)
	fmt.Printf("\n Folders for orgID: %s", orgID)
	folder.PrettyPrint(orgFolder)

	// example usage of get all child folders
	rootFolderString := "noble-vixen" // subject to change
	childFolders, err := folderDriver.GetAllChildFolders(orgID, rootFolderString)
	if err != nil {
		fmt.Println("\n")
		fmt.Println(err)
	} else {
		fmt.Println("\n\n Child folders of " + rootFolderString + ":")
		folder.PrettyPrint(childFolders)
	}

	// example use of move folder with sample json, note can use the example scenario json as readme reference
	sourceName := "nearby-secret"
	destinationName := "fast-watchmen"
	switchedFolders, err := folderDriver.MoveFolder(sourceName, destinationName)
	if err != nil {
		fmt.Println("\n")
		fmt.Println(err)
	} else {
		fmt.Println("Switched folders")
		folder.PrettyPrint(switchedFolders)
	}
}
