/**

 filename  : drive_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package gdrive

import (
	"strings"
	"bufio"
	"io"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveService struct {
	service *drive.Service
}

func getCredentials() ([]byte, error) {
	env := os.Getenv("GOOGLE_DRIVE_CREDENTIALS")

	if env != "" {
		return []byte(env), nil
	}

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return os.ReadFile(filepath.Join(wd, "configs", "google", "credentials.json"))
}

func getToken() (*oauth2.Token, error) {
	env := os.Getenv("GOOGLE_TOKEN")

	if env != "" {
		token := &oauth2.Token{}

		err := json.Unmarshal(
			[]byte(env),
			token,
		)

		return token, err
	}

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return tokenFromFile(filepath.Join(wd, "configs", "google", "token.json"))
}

func NewDriveService() (*DriveService, error) {

	ctx := context.Background()

	credentials, err := getCredentials()

	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(
		credentials,
		drive.DriveFileScope,
	)

	if err != nil {
		return nil, err
	}

	client, err := getClient(config)

	if err != nil {
		return nil, err
	}

	service, err := drive.NewService(
		ctx,
		option.WithHTTPClient(client),
	)

	if err != nil {
		return nil, err
	}

	return &DriveService{
		service: service,
	}, nil
}

func getClient(
	config *oauth2.Config,
) (*http.Client, error) {

	token, err := getToken()

	if err != nil {

		token, err = getTokenFromWeb(config)

		if err != nil {
			return nil, err
		}

		wd, err := os.Getwd()

		if err != nil {
			return nil, err
		}

		saveToken(
			filepath.Join(wd, "configs", "google", "token.json"),
			token,
		)
	}

	return config.Client(
		context.Background(),
		token,
	), nil
}

func getTokenFromWeb(
	config *oauth2.Config,
) (*oauth2.Token, error) {
 
	authURL := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
	)
 
	fmt.Println("Open this URL in your browser:")
	fmt.Println(authURL)
 
	// Bug 4 fix: use bufio.Reader so the full redirect URL (which may
	// contain spaces or special chars) is read correctly, and trim
	// whitespace so copy-paste with a trailing newline doesn't break
	// the exchange.
	fmt.Print("Enter code: ")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading auth code: %w", err)
	}
	code = strings.TrimSpace(code)
 
	token, err := config.Exchange(
		context.Background(),
		code,
	)
 
	if err != nil {
		return nil, err
	}
 
	return token, nil
}

func tokenFromFile(
	filePath string,
) (*oauth2.Token, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	token := &oauth2.Token{}

	err = json.NewDecoder(file).Decode(token)

	return token, err
}

func saveToken(
	filePath string,
	token *oauth2.Token,
) error {

	file, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

func (d *DriveService) UploadFile(
	filePath string,
) (*drive.File, error) {

	fileHandle, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer fileHandle.Close()

	folderId := os.Getenv("GOOGLE_DRIVE_BACKUP_FOLDER_ID")
	file := &drive.File{
		Name: filepath.Base(filePath),
	}

	// optional folder
	if folderId != "" {
		file.Parents = []string{
			folderId,
		}
	}

	response, err := d.service.Files.Create(
		file,
	).Media(fileHandle).Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}


func (d *DriveService) UpdateFile(
	fileId string,
	filePath string,
) (*drive.File, error) {

	fileHandle, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer fileHandle.Close()

	file := &drive.File{
		Name: filepath.Base(filePath),
	}

	response, err := d.service.Files.Update(
		fileId,
		file,
	).Media(fileHandle).Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DriveService) DownloadFile(
	fileId string,
	outputPath string,
) error {
	response, err := d.service.Files.Get(fileId).Download()

	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = os.MkdirAll(
		filepath.Dir(outputPath),
		os.ModePerm,
	)

	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(
		file,
		response.Body,
	)

	return err
}