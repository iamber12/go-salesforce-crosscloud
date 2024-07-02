package crosscloud

import (
	"fmt"
	"go-salesforce-crosscloud/gdrive"
	"go-salesforce-crosscloud/salesforce"
)

type GCPData struct {
	GCPServiceAccountCredentials string
	SFDomain                     string
	SFConsumerKey                string
	SFConsumerSecret             string
}

type CrossCloud struct {
	gDrive     *gdrive.GDrive
	salesforce *salesforce.Salesforce
	data       *GCPData
}

func InitMigrationToGDrive(gcpCredentials, sfDomain, sfConsumerKey, sfConsumerSecret string) (*CrossCloud, error) {
	gDrive, err := gdrive.InitializeGDrive(gcpCredentials)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Google Drive client: %w", err)
	}

	sf, err := salesforce.InitializeSalesforce(sfDomain, sfConsumerKey, sfConsumerSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Salesforce client: %w", err)
	}

	data := &GCPData{
		GCPServiceAccountCredentials: gcpCredentials,
		SFDomain:                     sfDomain,
		SFConsumerKey:                sfConsumerKey,
		SFConsumerSecret:             sfConsumerSecret,
	}
	return &CrossCloud{
		gDrive:     gDrive,
		salesforce: sf,
		data:       data,
	}, nil
}

func (cc *CrossCloud) MoveAllFilesToGoogleDrive(email, folderName string) error {
	cvDataList, err := cc.salesforce.GetAllFilesFromSalesforce()
	if err != nil {
		return fmt.Errorf("error getting files from Salesforce: %w", err)
	}

	if len(cvDataList) > 0 {
		folderId, err := cc.checkAndCreateFolder(folderName, email)
		if err != nil {
			return fmt.Errorf("error creating folder in Google Drive: %w", err)
		}

		for _, cvData := range cvDataList {
			err = cc.gDrive.UploadFileToGoogleDrive(cvData, folderId)
			if err != nil {
				return fmt.Errorf("error uploading file to Google Drive: %w", err)
			}
		}
	}

	return nil
}

func (cc *CrossCloud) MoveSObjectFilesToGDrive(sObject string, email string) error {
	folderVsCvDataMap, err := cc.salesforce.GetSObjectRelatedFilesFromSalesforce(sObject)
	if err != nil {
		return fmt.Errorf("error getting files from Salesforce: %w", err)
	}

	for folderName, filesToUpload := range folderVsCvDataMap {
		folderId, err := cc.checkAndCreateFolder(folderName, email)
		if err != nil {
			return fmt.Errorf("error creating folder in Google Drive: %w", err)
		}
		for _, cvData := range filesToUpload {
			err = cc.gDrive.UploadFileToGoogleDrive(cvData, folderId)
			if err != nil {
				return fmt.Errorf("error uploading file to Google Drive: %w", err)
			}
		}
	}
	return nil
}

func (cc *CrossCloud) checkAndCreateFolder(folderName, email string) (string, error) {
	var err error
	folderId, exists := cc.gDrive.FolderMap[folderName]
	if !exists {
		folderId, err = cc.gDrive.CreateNewFolder(folderName, email)
		if err != nil {
			return "", err
		}
		cc.gDrive.FolderMap[folderName] = folderId
	}

	return folderId, nil
}
