package folder_test

import (
	"errors"
	"testing"
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/stretchr/testify/assert"
	"github.com/gofrs/uuid"
)

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name string
		folders []folder.Folder
		sourceName string
		destinationName string
		wantFolders []folder.Folder
		wantError error
	} {
		{
			name: "Invalid same source and destination",
			folders: []folder.Folder{},
			sourceName: "A",
			destinationName: "A",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrSourceToItself),
		},
		{
			name: "Non-existent source folder",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			sourceName: "D",
			destinationName: "B",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrSourceNotExists),
		},
		{
			name: "Non-existent destination folder",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			sourceName: "A",
			destinationName: "D",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrDestNotExist),
		},
		{
			name: "Invalid move source to child destination",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			sourceName: "A",
			destinationName: "B",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrSourceToChild),
		},
		{
			name: "Invalid move to different organisation",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.Must(uuid.NewV4()), Paths: "A.B.C"},
			},
			sourceName: "A",
			destinationName: "C",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrFolderToDiffOrg),
		},
		{
			name: "Move folder with empty path",
			folders: []folder.Folder{
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A."},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "B"},
			},
			sourceName: "A",
			destinationName: "B",
			wantFolders: []folder.Folder{},
			wantError: errors.New(folder.ErrInvalidFilePath),
		},
		{
			name: "Move folder with nested children",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			sourceName: "B",
			destinationName: "A",
			wantFolders: []folder.Folder{
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantError: nil,
		},
		{
			name: "Source already child",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
			},
			sourceName: "B",
			destinationName: "A",
			wantFolders: []folder.Folder{
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			get, err := f.MoveFolder(tt.sourceName, tt.destinationName)
			if err != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantFolders, get)
		})
	}
}
