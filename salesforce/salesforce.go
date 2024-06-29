package salesforce

import (
	"fmt"
	"github.com/k-capehart/go-salesforce"
	"io"
	"log"
	"net/http"
	"strings"
)

type ContentVersion struct {
	Id            string `soql:"selectColumn,fieldName=Id" json:"Id"`
	Title         string `soql:"selectColumn,fieldName=Title" json:"Title"`
	VersionData   string `soql:"selectColumn,fieldName=VersionData" json:"VersionData"`
	FileExtension string `soql:"selectColumn,fieldName=FileExtension" json:"FileExtension"`
}

type CVData struct {
	cvData        []byte
	fileExtension string
}

type ContentDocumentLink struct {
	ContentDocumentId string       `json:"ContentDocument.LatestPublishedVersionId"`
	LinkedEntityId    string       `json:"LinkedEntityId"`
	LinkedEntity      LinkedEntity `json:"LinkedEntity"`
}

type LinkedEntity struct {
	Name string `json:"Name"`
}

type Salesforce struct {
	Client *salesforce.Salesforce
	Domain string
}

func InitializeSalesforce(domain, consumerKey, consumerSecret string) (*Salesforce, error) {
	sf, err := salesforce.Init(salesforce.Creds{
		Domain:         domain,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	})

	if err != nil {
		return nil, err
	}

	return &Salesforce{Client: sf, Domain: domain}, nil
}

func (sf *Salesforce) GetAllFilesFromSalesforce() ([]ContentVersion, error) {
	contentVersions := []ContentVersion{}
	contentVersionSoqlQuery := struct {
		SelectClause ContentVersion `soql:"selectClause,tableName=ContentVersion"`
	}{}

	err := sf.Client.QueryStruct(contentVersionSoqlQuery, &contentVersions)
	if err != nil {
		return nil, err
	}
	return contentVersions, nil
}

func (sf *Salesforce) queryContentDocumentLinkBySObject(sObject string) ([]ContentDocumentLink, error) {
	contentDocumentLinks := []ContentDocumentLink{}
	query := fmt.Sprintf(`SELECT ContentDocumentId, LinkedEntityId, LinkedEntity.Name FROM ContentDocumentLink WHERE LinkedEntityId IN (SELECT Id FROM %s)`, sObject)

	err := sf.Client.Query(query, &contentDocumentLinks)
	if err != nil {
		log.Panic(err)
	}

	return contentDocumentLinks, nil
}

func (sf *Salesforce) queryContentVersionByDocumentLink(contentDocumentIds []string) ([]ContentVersion, error) {
	contentVersions := []ContentVersion{}
	idsForQuery := "'" + strings.Join(contentDocumentIds, "', '") + "'"
	query := fmt.Sprintf("SELECT Id, Title, VersionData, FileExtension FROM ContentVersion WHERE ContentDocumentId IN (%s)", idsForQuery)

	err := sf.Client.Query(query, &contentVersions)
	if err != nil {
		log.Panic(err)
	}

	return contentVersions, nil
}

func (sf *Salesforce) GetSObjectRelatedFilesFromSalesforce(sObject string) (map[string][]CVData, error) {
	var contentDocumentIds []string
	var linkedEntityIds []string
	var folderVsCvDataMap = make(map[string][]CVData)
	var contentVersionVsFolderNameMap = make(map[string]string)

	contentDocumentLinks, err := sf.queryContentDocumentLinkBySObject(sObject)
	if err != nil {
		fmt.Printf("Error querying Salesforce: %v\n", err)
		return nil, err
	}

	for _, contentDocumentLink := range contentDocumentLinks {
		contentVersionVsFolderNameMap[contentDocumentLink.ContentDocumentId] = contentDocumentLink.LinkedEntity.Name
		contentDocumentIds = append(contentDocumentIds, contentDocumentLink.ContentDocumentId)
		linkedEntityIds = append(linkedEntityIds, contentDocumentLink.LinkedEntityId)
	}

	contentVersions, err := sf.queryContentVersionByDocumentLink(contentDocumentIds)
	if err != nil {
		fmt.Printf("Error querying Salesforce: %v\n", err)
		return nil, err
	}

	for _, contentVersion := range contentVersions {
		cvData, err := sf.DownloadFile(contentVersion.VersionData)
		if err != nil {
			return nil, err
		}

		folderName := contentVersionVsFolderNameMap[contentVersion.Id]
		newCVData := struct {
			cvData        []byte
			fileExtension string
		}{
			cvData:        cvData,
			fileExtension: contentVersion.FileExtension,
		}

		if folderVsCvDataMap[folderName] == nil {
			folderVsCvDataMap[folderName] = []CVData{newCVData}
		} else {
			folderVsCvDataMap[folderName] = append(folderVsCvDataMap[folderName], newCVData)
		}
	}

	return folderVsCvDataMap, nil
}

func (sf *Salesforce) DownloadFile(versionData string) ([]byte, error) {
	url := sf.Domain + versionData
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+sf.Client.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileData, err := io.ReadAll(resp.Body)
	return fileData, nil
}
