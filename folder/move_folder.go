package folder

import (
	"errors"
	"github.com/gofrs/uuid"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	if name == dst {
		return []Folder{}, errors.New(ErrSourceToItself)
	}

	folders := f.folders
	var nameFolder Folder // source folder
	var dstFolder Folder

	//Searching for origin IDs
	for _, f := range folders {
		if f.Name == name {
			nameFolder = f
		} else if f.Name == dst {
			dstFolder = f
		}
	}
	if nameFolder.OrgId == uuid.Nil {
		return []Folder{}, errors.New(ErrSourceNotExists) 
	}
	if dstFolder.OrgId == uuid.Nil {
		return []Folder{}, errors.New(ErrDestNotExist)
	}
	if nameFolder.OrgId != dstFolder.OrgId {
		return []Folder{}, errors.New(ErrFolderToDiffOrg)
	}

	childFolders, err := f.GetAllChildFolders(nameFolder.OrgId, name)
	if err != nil {
		return []Folder{}, err
	}

	// checking if destination is child of source
	for _, f := range childFolders {
		if f.Name == dst {
			return []Folder{}, errors.New(ErrSourceToChild)
		}
	}

	orgNamePathLen := len(nameFolder.Paths) // needed for path splitting
	newNamePath := dstFolder.Paths + "." + nameFolder.Name // new path prefix
	nameFolder.Paths = newNamePath

	// map used to save computation time for updating source + child folders
	updatingFolders := make(map[string]Folder) // map of name : Folder
	updatingFolders[name] = nameFolder
	for _, f := range childFolders {
		f.Paths = newNamePath + "." + f.Paths[orgNamePathLen + 1:] // (+1) due to extra '.'
		updatingFolders[f.Name] = f
	}

	for i := range folders {
		if _, exists := updatingFolders[folders[i].Name]; exists {
			folders[i].Paths = updatingFolders[folders[i].Name].Paths
		}
	}

	return folders, nil
}
