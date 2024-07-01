package gdrive

import (
	"context"
	"fmt"
	"go-salesforce-crosscloud/salesforce"
	"google.golang.org/api/option"
	"strings"

	drive "google.golang.org/api/drive/v3"
)

type GDrive struct {
	srv       *drive.Service
	FolderMap map[string]string
}

const SCOPE = drive.DriveScope

func InitializeGDrive(credentialFilePath string) (*GDrive, error) {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialFilePath), option.WithScopes(SCOPE))

	if err != nil {
		return nil, err
	}

	gDrive := &GDrive{srv: srv}

	gDrive.FolderMap, err = gDrive.QueryAllFolders()
	if err != nil {
		return nil, err
	}

	return gDrive, nil
}

func (gDrive *GDrive) QueryAllFolders() (map[string]string, error) {
	folderMap := make(map[string]string)
	query := "mimeType='application/vnd.google-apps.folder' and trashed=true"
	r, err := gDrive.srv.Files.List().Q(query).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	for _, folder := range r.Files {
		folderMap[folder.Name] = folder.Id
	}

	return folderMap, nil
}

func (gDrive *GDrive) CreateNewFolder(folderName string, email string) (string, error) {
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	folder, err := gDrive.srv.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder: %v", err)
	}

	if email != "" {
		permission := &drive.Permission{
			Type: "user",
			Role: "write",
			// Role:         "owner",
			EmailAddress: email,
		}

		_, err = gDrive.srv.Permissions.Create(folder.Id, permission).Do()
		if err != nil {
			return "", fmt.Errorf("unable to create permission: %v", err)
		}
	}

	return folder.Id, nil
}

func (gDrive *GDrive) UploadFileToGoogleDrive(cvData salesforce.CVData, folderId string) error {
	f := &drive.File{Name: cvData.Title + "." + cvData.FileExtension, Parents: []string{folderId}}
	_, err := gDrive.srv.Files.Create(f).Media(strings.NewReader(string(cvData.Data))).Do()
	if err != nil {
		return err
	}

	return nil
}
