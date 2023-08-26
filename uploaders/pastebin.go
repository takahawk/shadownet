package uploaders

import (
	"errors"

	"io/ioutil"

	"net/http"
	"net/url"

	"regexp"
	"strings"
)

const PastebinUploaderName = "pastebin"

const MaximumPastebinExpireTime = "1Y"
const PastebinUploadApiOption = "paste"
const (
	PastebinPrivacyPublic = "0"
	PastebinPrivacyUnlisted = "1"
	PastebinPrivacyPrivate = "2"
)

var successfulResponsePattern = regexp.MustCompile(`https://pastebin.com/(.+)$`)
	

type pastebinUploader struct {
	apiKey string
}

// KIM:
// 1) Maximum storage time is one year
// 2) Maximum number of unlisted pastes for free account are 25
func NewPastebinUploader(apiKey string) Uploader {
	return &pastebinUploader{
		apiKey: apiKey,
	}
}

func NewPastebinUploaderWithParams(params... string) (Uploader, error) {
	if len(params) != 1 {
		return nil, errors.New("there should be 1 parameters: pastebin developer key")
	}
	apiKey := params[0]

	return &pastebinUploader{apiKey}, nil
}

func (pu *pastebinUploader) Name() string {
	return PastebinUploaderName
}

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

	return groups[1], nil
}