package uploaders

import (
	"errors"
	"fmt"

	"io/ioutil"

	"net/http"
	"net/url"

	"regexp"
	"strings"

	"github.com/takahawk/shadownet/logger"
)

// PastebinUploaderName is pastebin uploader component name
const PastebinUploaderName = "pastebin"

// TODO: mb should use N?
// MaximumPastebinExpireTime is maximum value for `api_paste_expire_date` and
// it means 1 year. Other possible values are: N (never), 10M (10 minutes),
// 1H (1 hour), 1D (1 day), 1W (1 week), 2W (2 week), 1M (1 month),
// 6M (6 months)
const MaximumPastebinExpireTime = "1Y"

// PastebinUploadApiOption is value for `api_option` parameter that should be
// set to this for upload operation
const PastebinUploadApiOption = "paste"

const (
	// PastebinPrivacyPublic is parameter for `api_paste_private` used to make
	// paste public (available for anyone and listed in search results)
	PastebinPrivacyPublic = "0"
	// PastebinPrivacyUnlisted is parameter for `api_paste_private` used to make
	// paste unlisted (available for anyone but not listed in search results)
	PastebinPrivacyUnlisted = "1"
	// PastebinPrivacyPrivate is parameter for `api_paste_private` used to make
	// paste private (accessible only when logged in to corresponding account)
	PastebinPrivacyPrivate = "2"
)

// PastebinRawPrefix is prefix for URL used to get saved paste in raw
// (e.g. https://pastebin.com/raw/y1FKvrXe)
const PastebinRawPrefix = "https://pastebin.com/raw"

var successfulResponsePattern = regexp.MustCompile(`https://pastebin.com/(.+)$`)

type pastebinUploader struct {
	logger logger.Logger
	apiKey string
}

// NewPastebinUploader returns uploader that uploads data to Pastebin for a
// given API key that will be used to make upload requests
func NewPastebinUploader(logger logger.Logger, apiKey string) Uploader {
	return &pastebinUploader{
		apiKey: apiKey,
	}
}

// NewPastebinUploaderWithParams returns uploader for a given params. It
// does expect single param that is API key. It exists only for convenience
// doing effectively the same as NewPastebinUploader
func NewPastebinUploaderWithParams(logger logger.Logger, params ...[]byte) (Uploader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be 1 parameters: pastebin developer key")
	}
	apiKey := string(params[0])

	return &pastebinUploader{
		logger: logger,
		apiKey: apiKey,
	}, nil
}

// Name returns pastebin uploader name. It is always PastebinUploaderName
func (pu *pastebinUploader) Name() string {
	return PastebinUploaderName
}

// Params returns API key packed into byte array
func (pu *pastebinUploader) Params() [][]byte {
	return [][]byte{[]byte(pu.apiKey)}
}

// Upload saves data in byte array as a paste on Pastebin
func (pu *pastebinUploader) Upload(content []byte) (id string, err error) {
	// TODO: should check if this is possible to upload binary data
	// TODO: implement
	apiUrl := "https://pastebin.com"
	resource := "/api/api_post.php"
	data := url.Values{}
	data.Set("api_dev_key", pu.apiKey)
	data.Set("api_option", PastebinUploadApiOption)
	data.Set("api_paste_expire_date", MaximumPastebinExpireTime)
	data.Set("api_paste_private", PastebinPrivacyPublic)
	data.Set("api_paste_code", string(content))

	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		return "", err
	}
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	groups := successfulResponsePattern.FindStringSubmatch(string(body))
	if len(groups) != 2 {
		return "", errors.New("failed to capture id")
	}

	id = fmt.Sprintf("%s/%s", PastebinRawPrefix, groups[1])

	return id, nil
}
