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

/* Validates previous folders have previously been seen */
func ValidatePathStructure(path string, seen map[string]int) error {
	splitPaths := strings.Split(path, ".") // Expects current folder to be a child to a previous folder
	if len(splitPaths) < 2 {
		return errors.New("Error: invalid file path structure " + path)
	}
	/* all previous files must be seen in order for the current path to be valid because we have 
	sorted the folders, the previous folder must have already been seen */
	if _, exists := seen[splitPaths[len(splitPaths) - 2]]; exists {
		return nil
	}
	return errors.New("Error: path contains unseen folder " + path + " for " + splitPaths[len(splitPaths) - 2]) 
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

				err := ValidatePathStructure(f.Paths, seen)
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
