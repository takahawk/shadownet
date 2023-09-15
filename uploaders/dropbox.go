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

const DropboxUploaderName = "dropbox"

const DropboxApiUrlUpload = "https://content.dropboxapi.com/2/files/upload"
const DropboxApiUrlCreateSharedLink = "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

const (
	DropboxUploadArgModeAdd       = "add"
	DropboxUploadArgModeOverwrite = "overwrite"
	DropboxUploadArgModeUpdate    = "update"
)

const (
	DropboxCreateSharedLinkAccessViewer  = "viewer"
	DropboxCreateSharedLinkAccessEditor  = "editor"
	DropboxCreateSharedLinkAccessMax     = "max"
	DropboxCreateSharedLinkAccessDefault = "default"
)

const (
	DropboxCreateSharedLinkAudiencePublic = "public"
	DropboxCreateSharedLinkAudienceTeam   = "team"
	DropboxCreateSharedLinkAudienceNoOne  = "no_one"
)

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

func NewDropboxUploader(logger logger.Logger, accessToken string) Uploader {
	return &dropboxUploader{
		accessToken: accessToken,
		logger:      logger,
	}
}

func NewDropboxUploaderWithParams(logger logger.Logger, params ...[]byte) (Uploader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be exactly 1 parameter: Dropbox access token")
	}
	return NewDropboxUploader(logger, string(params[0])), nil
}

func (du *dropboxUploader) Name() string {
	return DropboxUploaderName
}

func (du *dropboxUploader) Params() [][]byte {
	return [][]byte{[]byte(du.accessToken)}
}

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
			Access:        DropboxCreateSharedLinkAccessMax,
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
