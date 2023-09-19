package uploaders

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/takahawk/shadownet/logger"
)

// DropboxUploaderName is dropbox uploader component name
const DropboxUploaderName = "dropbox"

// DropboxApiUrlUpload is URL to send POST upload requests to dropbox
const DropboxApiUrlUpload = "https://content.dropboxapi.com/2/files/upload"

// DropboxApiUrlCreateSharedLink is URL to send POST request to create shared
// link for already uploaded file
const DropboxApiUrlCreateSharedLink = "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

const (
	// DropboxUploadArgModeAdd is write mode that in case of conflict creates
	// file with another name (i.e. "filename(2).txt")
	DropboxUploadArgModeAdd = "add"
	// DropboxUploadArgModeOverwrite is write mode that in case of conflict
	// just overwrites
	DropboxUploadArgModeOverwrite = "overwrite"
	// DropboxUploadArgModeUpdate is write mode that overwrites file, but
	// only if it's "rev" matches the existing file's "rev" (see more in
	// Dropbox HTTP API docs)
	DropboxUploadArgModeUpdate = "update"
)

const (
	// DropboxCreateSharedLinkAccessViewer is access level that allows users
	// who has link to view and comment on the content
	DropboxCreateSharedLinkAccessViewer = "viewer"
	// DropboxCreateSharedLinkAccessEditor is access level that allows users
	// who has link to view, comment and also edit the content (not all files
	// support edit links)
	DropboxCreateSharedLinkAccessEditor = "editor"
	// DropboxCreateSharedLinkAccessMax is  sets access level to maximum
	// allowed level for this URL
	DropboxCreateSharedLinkAccessMax = "max"
	// DropboxCreateSharedLinkAccessDefault is set access level to the
	// default user has set
	DropboxCreateSharedLinkAccessDefault = "default"
)

const (
	// DropboxCreateSharedLinkAudiencePublic makes link accessible to anyone
	DropboxCreateSharedLinkAudiencePublic = "public"
	// DropboxCreateSharedLinkAudienceTeam makes link accessible only by team
	// members
	DropboxCreateSharedLinkAudienceTeam = "team"
	// DropboxCreateSharedLinkAudienceNoOne make link to not grant any
	// additional access rights. Link only points to content and can
	// be only successfully used by user who already have some access
	// to a file
	DropboxCreateSharedLinkAudienceNoOne = "no_one"
)

// RandomFilenameBytes is number used to generate random name for a data to
// be uploaded
const RandomFilenameBytes = 16

type dropboxUploader struct {
	accessToken string
	logger      logger.Logger
}

type dropboxUploadApiArg struct {
	Autorename     bool   `json:"autorename"`
	Mode           string `json:"mode"`
	Mute           bool   `json:"mute"`
	Path           string `json:"path"`
	StrictConflict bool   `json:"strict_conflict"`
}

type dropboxCreateSharedLinkRequestBody struct {
	Path     string `json:"path"`
	Settings struct {
		Access        string `json:"access"`
		AllowDownload bool   `json:"allow_download"`
		Audience      string `json:"audience"`
	} `json:"settings"`
}

type dropboxCreateSharedLinkResponseBody struct {
	Url string `json:"url"`
}

// NewDropboxUploader return new dropbox uploader that will be using given
// access token
func NewDropboxUploader(logger logger.Logger, accessToken string) Uploader {
	return &dropboxUploader{
		accessToken: accessToken,
		logger:      logger,
	}
}

// NewDropboxUploaderWithParams is convinience method that calls NewDropboxUploader
// but with access token packed into byte slice
func NewDropboxUploaderWithParams(logger logger.Logger, params ...[]byte) (Uploader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be exactly 1 parameter: Dropbox access token")
	}
	return NewDropboxUploader(logger, string(params[0])), nil
}

// Name returns Dropbox uploader component name. It is always DropboxUploaderName
func (du *dropboxUploader) Name() string {
	return DropboxUploaderName
}

// Params returns access token packed into byte slice
func (du *dropboxUploader) Params() [][]byte {
	return [][]byte{[]byte(du.accessToken)}
}

// Upload uploads given file to dropbox and returns shared link to it
func (du *dropboxUploader) Upload(data []byte) (id string, err error) {
	du.logger.Info("Uploading data to Dropbox...")
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, DropboxApiUrlUpload, bytes.NewReader(data))
	if err != nil {
		du.logger.Errorf("%+v", err)
		return "", err
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", du.accessToken))

	id = generateRandomFilename()
	args := dropboxUploadApiArg{
		Autorename:     false,
		Mode:           DropboxUploadArgModeAdd,
		Mute:           true,
		Path:           fmt.Sprintf("/%s", id),
		StrictConflict: true,
	}
	apiArgs, err := json.Marshal(args)
	if err != nil {
		du.logger.Errorf("%+v", err)
		return "", err
	}

	r.Header.Set("Dropbox-API-Arg", string(apiArgs))
	r.Header.Set("Content-Type", "application/octet-stream")

	rsp, err := client.Do(r)
	if err != nil {
		du.logger.Errorf("%+v", err)
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		du.logger.Errorf("Request failed with status code: %d", rsp.StatusCode)
		body, err := io.ReadAll(rsp.Body)
		if err != nil {
			du.logger.Errorf("Error reading response body: %s", string(body))
		}
		du.logger.Errorf(string(body))
		return "", errors.New("request failed")
	}
	du.logger.Infof("Success uploading data to Dropbox. ID: %s", id)

	// TODO: refactor, move all common dropbox HTTP request logic to separate function
	body := dropboxCreateSharedLinkRequestBody{
		Path: fmt.Sprintf("/%s", id),
		Settings: struct {
			Access        string `json:"access"`
			AllowDownload bool   `json:"allow_download"`
			Audience      string `json:"audience"`
		}{
			Access:        DropboxCreateSharedLinkAccessViewer,
			AllowDownload: true,
			Audience:      DropboxCreateSharedLinkAudiencePublic,
		},
	}
	bodyJson, err := json.Marshal(body)
	du.logger.Infof("Creating shared link...")
	r, err = http.NewRequest(http.MethodPost, DropboxApiUrlCreateSharedLink, bytes.NewReader(bodyJson))
	if err != nil {
		du.logger.Errorf("%+v", err)
		return "", err
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", du.accessToken))
	r.Header.Set("Content-Type", "application/json")

	rsp, err = client.Do(r)
	if err != nil {
		du.logger.Errorf("%+v", err)
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		du.logger.Errorf("Request failed with status code: %d", rsp.StatusCode)
		body, err := io.ReadAll(rsp.Body)
		if err != nil {
			du.logger.Error("Error reading response body")
			return "", errors.New("request failed")
		}
		du.logger.Errorf(string(body))
		return "", errors.New("request failed")
	}

	rspBodyJson, err := io.ReadAll(rsp.Body)
	if err != nil {
		du.logger.Error("Error reading response body")
		return "", errors.New("request failed")
	}
	var rspBody dropboxCreateSharedLinkResponseBody
	err = json.Unmarshal(rspBodyJson, &rspBody)
	if err != nil {
		du.logger.Errorf("Error unmarshalling response body: %s", string(rspBodyJson))
		return "", errors.New("request failed")
	}

	// substitute dl=0 to dl=1 to get direct download link
	id = rspBody.Url[:len(rspBody.Url)-1] + "1"
	du.logger.Infof("Success creating shared link. Web URL: %s", id)

	return
}

// TODO: mb create separate utils package for such things
func generateRandomFilename() string {
	randBytes := make([]byte, RandomFilenameBytes)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}
