package folder_test

import (
	"errors"
	"testing"
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		// TODO: your tests here
		{
			name: "Only default org ID",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			folders: folder.GetSampleData(),
			want: folder.GetSampleDefaultOrgIDOnlyData(), // Manually extracted out of the sample.json
		},
		{
			name: "Not existing OrgID",
			orgID: uuid.Must(uuid.NewV4()),
			folders: []folder.Folder{},
			want: []folder.Folder{},
		},
		{
			name: "Existing OrgID with empty folder",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			folders: []folder.Folder{},
			want: []folder.Folder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			get := f.GetFoldersByOrgID(tt.orgID)
			assert.Equal(t, tt.want, get)
		})
	}
}

func Test_folder_ValidateFilePath(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name string
		path string
		want bool
	} {
		{
			name: "Root",
			path: "A",
			want: true,
		},
		{
			name: "Root & Child",
			path: "A.B",
			want: true,
		},
		{
			name: "Root & Child Folder & Child's Folder",
			path: "A.B.C",
			want: true,
		},
		{
			name: "Invalid missing ending folder",
			path: "A.",
			want: false,
		},
		{
			name: "Invalid missing ending folder with child",
			path: "A.B.",
			want: false,
		},
		{
			name: "Numbered folder name",
			path: "1",
			want: true,
		},
		{
			name: "Numbered folder name & child",
			path: "1.2",
			want: true,
		},
		{
			name: "Empty path",
			path: "",
			want: false,
		},
		{
			name: "Just dot",
			path: ".",
			want: false,
		},
		{
			name: "dot then folder",
			path: ".A",
			want: false,
		},
		{
			name: "Path folder with numbers & letters",
			path: "A1",
			want: true,
		},
		{
			name: "Path folder with numbers & letters with child",
			path: "A1.B2",
			want: true,
		},
		{
			name: "Path folder a special character",
			path: "!",
			want: false,
		},
		{
			name: "Path folder with special characters",
			path: "A.%",
			want: false,
		},
		{
			name: "Path folder with consecutive separators",
			path: "A..B",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			get := folder.ValidateFilePath(tt.path)
			assert.Equal(t, tt.want, get)
		})
	}
}

func Test_folder_ValidateFolderEndOfPath(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name string
		folder folder.Folder
		want bool
	} {
		{
			name: "Root folder",
			folder: folder.GetSampleData()[0],
			want: true,
		},
		{
			name: "Child folder",
			folder: folder.GetSampleData()[1],
			want: true,
		},
		{
			name: "Path name mismatch",
			folder: folder.Folder{
				Name: "C",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "A.C.D",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, folder.ValidateFolderEndOfPath(tt.folder))
		})
	}
}

func Test_folder_ValidateChildPathStructure(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name string
		path string
		seen map[string]int
		want error
	} {
		{
			name: "Valid folder path with seen parent",
			path: "A.B.C",
			seen: map[string]int{"B": 1},
			want: nil,
		},
		{
			name: "Invalid folder path with unseen parent",
			path: "A.B.C",
			seen: map[string]int{"A": 1},
			want: errors.New(folder.ErrUnseenFolder + " A.B.C for B"),
		},
		{
			name: "Invalid short folder path",
			path: "A",
			seen: map[string]int{},
			want: errors.New(folder.ErrInvalidFilePathStructure + " A"),
		},
		{
			name: "Edge Case: Empty path",
			path: "",
			seen: map[string]int{},
			want: errors.New(folder.ErrInvalidFilePathStructure + " "),
		},
		{
			name: "Edge case: Missing child folder path",
			path: "A.",
			seen: map[string]int{},
			want: errors.New(folder.ErrInvalidFilePathStructure + " A."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T)  {
			assert.Equal(t, tt.want, folder.ValidateChildPathStructure(tt.path, tt.seen))
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name string
		orgID uuid.UUID
		rootFolderName string
		folders []folder.Folder
		wantFolders []folder.Folder
		wantErr error
	} {
		{
			name: "Valid direct child folder",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			rootFolderName: "A",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
			},
			wantFolders: []folder.Folder {
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
			},
			wantErr: nil,
		},
		{
			name: "Valid sub folder",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			rootFolderName: "A",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantFolders: []folder.Folder {
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantErr: nil,
		},
		{
			name: "Valid root sub folder",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			rootFolderName: "B",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantFolders: []folder.Folder {
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantErr: nil,
		},
		{
			name: "Valid leaf folder",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			rootFolderName: "C",
			folders: []folder.Folder {
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
				{Name: "C", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B.C"},
			},
			wantFolders: []folder.Folder {},
			wantErr: nil,
		},
		{
			name: "Invalid nil orgID",
			orgID: uuid.Nil,
			rootFolderName: "A",
			folders: []folder.Folder{},
			wantErr: errors.New(folder.ErrInvalidOrgID),
		},
		{
			name: "No folders in the organization",
			orgID: uuid.Must(uuid.NewV4()),
			rootFolderName: "A",
			folders: []folder.Folder{},
			wantErr: errors.New(folder.ErrFolderNotExistsOrg),
		},
		{
			name: "Root folder does not exist",
			orgID: uuid.FromStringOrNil(folder.DefaultOrgID),
			rootFolderName: "Empty", 
			folders: []folder.Folder{
				{Name: "A", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A"},
				{Name: "B", OrgId: uuid.FromStringOrNil(folder.DefaultOrgID), Paths: "A.B"},
			},
			wantErr: errors.New(folder.ErrFolderNotExist),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			get, err := f.GetAllChildFolders(tt.orgID, tt.rootFolderName)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantFolders, get)
		})
	}
}