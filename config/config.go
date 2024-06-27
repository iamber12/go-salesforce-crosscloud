package config

import (
	"context"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
)

const SCOPE = drive.DriveScope

func InitializeDriveService(serviceAccountCredentials string, serviceAccountPath string) (*drive.Service, error) {
	ctx := context.Background()

	credentialsBytes := []byte(serviceAccountCredentials)

	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialsBytes), option.WithScopes(SCOPE))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}
	return srv, nil
}
