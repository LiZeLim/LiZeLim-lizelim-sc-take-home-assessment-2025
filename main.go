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
	res = folder.GetSampleData()
	folderDriver = folder.NewDriver(res)
	orgID = uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")

	rootFolderString := "clear-arclight" // subject to change
	arcLightChildren, err := folderDriver.GetAllChildFolders(orgID, rootFolderString)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\n\n Child folders of " + rootFolderString + ":")
		folder.PrettyPrint(arcLightChildren)
	}
}
