package downloaders

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const PastebinRawPrefix = "https://pastebin.com/raw"
const PastebinPostUrl = "https://pastebin.com/api/api_post.php"

type pastebinDownloader struct {
	pasteID string
}
const PastebinDownloaderName = "pastebin"

func NewPastebinDownloader(pasteID string) Downloader {
	return &pastebinDownloader{
		pasteID: pasteID,
	}
}

func NewPastebinDownloaderWithParams(params... []byte) (Downloader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be only 1 param: paste ID")
	}
	return NewPastebinDownloader(string(params[0])), nil
}

func (pd *pastebinDownloader) Name() string {
	return PastebinDownloaderName
}

func (pd *pastebinDownloader) Params() [][]byte {
	return [][]byte{ []byte(pd.pasteID) }
}

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
