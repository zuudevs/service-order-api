/**

 filename  : drive_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveService struct {
	service *drive.Service
}

func NewDriveService() (*DriveService, error) {

	ctx := context.Background()

	credentials, err := os.ReadFile(
		"configs/google/credentials.json",
	)

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

	token, err := tokenFromFile(
		"configs/google/token.json",
	)

	if err != nil {

		token, err = getTokenFromWeb(config)

		if err != nil {
			return nil, err
		}

		saveToken(
			"configs/google/token.json",
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

	fmt.Println(
		"Open this URL in your browser:",
	)

	fmt.Println(authURL)

	fmt.Print("Enter code: ")

	var code string

	fmt.Scan(&code)

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

	file := &drive.File{
		Name: filePath,
	}

	response, err := d.service.Files.Create(
		file,
	).Media(fileHandle).Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}