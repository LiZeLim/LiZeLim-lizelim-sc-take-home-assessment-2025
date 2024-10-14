package folder

import (
	//"fmt"
	"errors"
	"sort"
	"strings"
	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res
}

/* Checks if the end of path matches the folder name */
func ValidateFolderEndOfPath(folder Folder) bool {
	pathLength := len(folder.Paths)
	nameLength := len(folder.Name)
	return folder.Paths[pathLength - nameLength:] == folder.Name
}

/* Validates if all previous folders up to the root folder have previously been seen */
func validatePathStructure(path string, rootFoldername string, seen map[string]int) error {
	splitPaths := strings.Split(path, ".")
	if len(splitPaths) < 2 {
		return errors.New("Error: invalid file path structure " + path)
	}
	for i := len(splitPaths) - 2; i >= 0; i-- { // all previous files must be seen in order for the current path to be valid
		if _, exists := seen[splitPaths[i]]; exists {
			if splitPaths[i] == rootFoldername {
				break
			}
			continue
		} else {
			return errors.New("Error: path contains unseen folder " + path + " for " + splitPaths[i]) 
		}
	}
	return nil
}

/* 
	The main premise of the algorithm is to sort the list of folders by their path name.
	By lexicographically sorting the paths, we can ensure that they are ordered correctly by their
	hierarchy. 
	*/
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	if orgID.IsNil() {
		return nil, errors.New("Error: Invalid orgID")
	}

	sameOriginFolders := f.GetFoldersByOrgID(orgID)

	if len(sameOriginFolders) == 0 {
		return nil, errors.New("Error: Folder does not exist in the specified organization")
	}


	/* sorting the folders of same origin to ensure child folders are in correct 
	ordering, also allow for prev folder checking. */ 
	sort.SliceStable(sameOriginFolders, func(i, j int) bool {
		return sameOriginFolders[i].Paths < sameOriginFolders[j].Paths
	})


	var rootFolder *Folder
	seen := make(map[string]int)

	res := []Folder{}
	for i := range sameOriginFolders {
		f := &sameOriginFolders[i]

		// finding root folder
		if rootFolder == nil && f.Name == name {
			rootFolder = f
			seen[f.Name] = 1
			continue // continue, root folder is not child folder
		}

		if rootFolder != nil { 
			if len(f.Paths) > len(rootFolder.Paths) && f.Paths[:len(rootFolder.Paths)] == rootFolder.Paths {
				if !ValidateFolderEndOfPath(*f) {
					return nil, errors.New("Error: Folder name doesn't match end of path " + f.Paths)
				}

				err := validatePathStructure(f.Paths, rootFolder.Paths, seen)
				if err != nil {
					return nil, err
				}

				res = append(res, *f)
				seen[f.Name] = 1
			} else {
				break //early exit since list is sorted, no other child will exist if previous conditions fail
			}
		}
	}

	if rootFolder == nil {
		return nil, errors.New("Error: Folder does not exist")
	}

	return res, nil
}
