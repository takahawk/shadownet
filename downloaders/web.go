package downloaders

import (
	"errors"
	"io"
	"net/http"

	"github.com/takahawk/shadownet/logger"
)

type webDownloader struct {
	logger logger.Logger
	url    string
}

// WebDownloaderName is web downloader component name
const WebDownloaderName = "web"

// NewWebDownloader returns downloader that downloads the data by the given
// URL using HTTP request
func NewWebDownloader(logger logger.Logger, url string) Downloader {
	return &webDownloader{
		logger: logger,
		url:    url,
	}
}

// NewWebDownloaderWithParams returns downloader for a given params. It
// does expect single param that is url. It exists only for convenience
// doing effectively the same as NewWebDownloader
func NewWebDownloaderWithParams(logger logger.Logger, params ...[]byte) (Downloader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be only 1 param: url")
	}
	return NewWebDownloader(logger, string(params[0])), nil
}

// Name returns web downloader name. It is always WebDownloaderName
func (wd *webDownloader) Name() string {
	return WebDownloaderName
}

// Params returns url packed into byte array
func (wd *webDownloader) Params() [][]byte {
	return [][]byte{[]byte(wd.url)}
}

// Download returns downloaded data in a byte array
func (wd *webDownloader) Download() ([]byte, error) {
	wd.logger.Infof("Downloading data from URL: %s", wd.url)
	res, err := http.Get(wd.url)
	if err != nil {
		wd.logger.Errorf("Error downloading data: %+v", err)
		// TODO: error handling (wrap etc.)?
		return nil, err
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		wd.logger.Errorf("Error reading response: %+v", err)
		return nil, err
	}

	wd.logger.Infof("Successfully downloaded")
	return content, nil
}
