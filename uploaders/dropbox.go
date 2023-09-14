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

const (
	DropboxUploadArgModeAdd       = "add"
	DropboxUploadArgModeOverwrite = "overwrite"
	DropboxUploadArgModeUpdate    = "update"
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

func NewDropboxUploader(accessToken string, logger logger.Logger) Uploader {
	return &dropboxUploader{
		accessToken: accessToken,
		logger:      logger,
	}
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

	return
}

// TODO: mb create separate utils package for such things
func generateRandomFilename() string {
	randBytes := make([]byte, RandomFilenameBytes)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}
