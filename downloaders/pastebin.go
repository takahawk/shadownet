package downloaders

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/takahawk/shadownet/logger"
)

// PastebinRawPrefix is prefix for URL used to get saved paste in raw
// (e.g. https://pastebin.com/raw/y1FKvrXe)
const PastebinRawPrefix = "https://pastebin.com/raw"

// PastebinPostUrl is URL to post new paste to Pastebin
const PastebinPostUrl = "https://pastebin.com/api/api_post.php"

// TODO: make generic Web downloader instead
type pastebinDownloader struct {
	logger  logger.Logger
	pasteID string
}

// PastebinDownloaderName is pastebin downloader component name
const PastebinDownloaderName = "pastebin"

// NewPastebinDownloader returns downloader for a given paste ID. Paste ID is
// the last part in the URL used to identify paste (e.g. y1FKvrXe in
// https://pastebin.com/raw/y1FKvrXe)
func NewPastebinDownloader(logger logger.Logger, pasteID string) Downloader {
	return &pastebinDownloader{
		logger:  logger,
		pasteID: pasteID,
	}
}

// NewPastebinDownloaderWithParams returns downloader for a given params. It
// does expect single param that is paste ID. It exists only for convenience
// doing effectively the same as NewPastebinDownloader
func NewPastebinDownloaderWithParams(logger logger.Logger, params ...[]byte) (Downloader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be only 1 param: paste ID")
	}
	return NewPastebinDownloader(logger, string(params[0])), nil
}

// Name returns pastebin downloader name. It is always PastebinDownloaderName
func (pd *pastebinDownloader) Name() string {
	return PastebinDownloaderName
}

// Params returns paste ID packed into byte array
func (pd *pastebinDownloader) Params() [][]byte {
	return [][]byte{[]byte(pd.pasteID)}
}

// Download returns downloaded paste in a byte array
func (pd *pastebinDownloader) Download() ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", PastebinRawPrefix, pd.pasteID))
	if err != nil {
		// TODO: error handling (wrap etc.)?
		return nil, err
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return content, nil
}
