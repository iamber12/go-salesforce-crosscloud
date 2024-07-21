# CrossCloud File Migration

This project provides a Go package for migrating files between different cloud services. Currently, it supports migrating files from Salesforce to Google Drive. It supports two modes: migrating all files or migrating files by specific objects. The package handles the creation of corresponding directories in Google Drive and uploads the files into those directories.

**Note: This project is a work in progress. Future updates will include support for additional cloud services such as SharePoint. We are actively looking for contributors to help us expand and improve this project.**

## Features

- **Move All Files**: Migrates all files from Salesforce to a specified Google Drive folder.
- **Move Files by Object**: Creates a directory for each Salesforce record and uploads the corresponding files into those directories.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/iamber12/go-salesforce-crosscloud
    cd go-salesforce-crosscloud
    ```

2. Install the required dependencies:
    ```bash
    go mod download
    ```

## Usage

The main package is `crosscloud`. Users do not need to interact with the `gdrive` and `salesforce` packages directly. Here is an example of how to use the package:

### Example

```go
package main

import "go-salesforce-crosscloud/crosscloud"

func main() {
    email := "your-email@example.com"
    folderName := "Salesforce Files"
    googleCredentialPath := "path/to/your/credentials.json"
    salesforceDomain := "your-salesforce-domain"
    salesforceConsumerKey := "your-salesforce-consumer-key"
    salesforceConsumerSecret := "your-salesforce-consumer-secret"

    cc, err := crosscloud.InitMigrationToGDrive(credentialPath, domain, consumerKey, consumerSecret)
    if err != nil {
        panic(err)
        return
    }

    err = cc.MoveAllFilesToGoogleDrive(email, folderName)
    // Uncomment the line below to move files by Salesforce object
    // err = cc.MoveSObjectFilesToGDrive("Case", email)
    if err != nil {
        panic(err)
        return
    }
}
```

## Configuration

### Google Drive Setup

1. Obtain Google Drive API service account credentials:
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project or select an existing project.
   - Enable the Google Drive API for your project.
   - Create a service account and download the JSON key file containing the credentials.
   - Place the JSON key file in your project directory.

### Salesforce Setup

1. Obtain Salesforce API credentials:
   - Log in to your Salesforce account.
   - Go to `Setup` > `Apps` > `App Manager`.
   - Create a new Connected App.
   - Enable OAuth settings and specify the callback URL.
   - Save the Consumer Key and Consumer Secret.

## Methods

### InitMigrationToGDrive

Initializes the migration process by setting up the Google Drive and Salesforce clients.

**Parameters:**
- `gcpCredentials`: Path to the Google Drive API credentials file.
- `sfDomain`: Salesforce domain.
- `sfConsumerKey`: Salesforce Consumer Key.
- `sfConsumerSecret`: Salesforce Consumer Secret.

**Returns:**
- A `CrossCloud` instance.

### MoveAllFilesToGoogleDrive

Moves all files from Salesforce to Google Drive.

**Parameters:**
- `email`: (Optional) Email address for sharing the folder. If the user does not have direct access to the service account, permissions will be manually granted to this email address.
- `folderName`: Name of the Google Drive folder.

**Returns:**
- `error`: An error if any occurred during the process.

### MoveSObjectFilesToGDrive

Moves files related to a specific Salesforce object to Google Drive.

**Parameters:**
- `sObject`: Salesforce object name (e.g., "Case").
- `email`: Email address for sharing the folder.

**Returns:**
- `error`: An error if any occurred during the process.

## Contribution Guidelines

We welcome contributions to this project! Do checkout the issues if you're interested. To contribute, follow these steps:

1. **Fork the Repository**
   - Go to the [https://github.com/iamber12/go-salesforce-crosscloud](https://github.com/iamber12/go-salesforce-crosscloud).
   - Click the "Fork" button in the upper right corner.

2. **Clone the Forked Repository**
   - Clone the forked repository to your local machine:
     ```bash
     git clone https://github.com/yourusername/go-salesforce-crosscloud.git
     cd go-salesforce-crosscloud
     ```

3. **Create a New Branch**
   - Create a new branch for your feature or bug fix:
     ```bash
     git checkout -b feature-name
     ```

4. **Make Changes**
   - Implement your changes in the new branch.
   - Follow the project's coding standards and best practices.

5. **Add Tests**
   - Add tests to ensure the reliability and correctness of your changes.

6. **Commit Changes**
   - Commit your changes with a descriptive commit message:
     ```bash
     git add .
     git commit -m "Description of changes"
     ```

7. **Push Changes**
   - Push your changes to your forked repository:
     ```bash
     git push origin feature-name
     ```

8. **Create a Pull Request**
   - Go to the original repository on GitHub.
   - Click the "Pull Requests" tab.
   - Click the "New Pull Request" button.
   - Select your branch from the "compare" dropdown.
   - Provide a detailed description of your changes and submit the pull request.

## Next Steps
- Create a new package for SharePoint integration, similar to the existing `gdrive` package. This package will handle authentication and file operations with SharePoint.
- Add multithreading to improve file download and upload functionality.
- Refactor the package structure so that only the `crosscloud` package is accessible to users.
- Create a Go package from this project to simplify installation and usage.
- Add tests to ensure the reliability and correctness of the migration functionality.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/iamber12/go-salesforce-crosscloud/blob/main/LICENSE) file for details.

## Acknowledgements

This project uses the following libraries:
- [googleapis/drive](https://pkg.go.dev/google.golang.org/api/drive/v3)
- [k-capehart/go-salesforce](https://github.com/k-capehart/go-salesforce)
