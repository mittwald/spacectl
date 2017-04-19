package lowlevel

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/mittwald/spacectl/service/auth"
	"bytes"
	"time"
)

type SpacesLowlevelClient struct {
	token string
	endpoint string
	version string

	client *http.Client
}

func NewSpacesLowlevelClient(token string, endpoint string) (*SpacesLowlevelClient, error) {
	client := &http.Client{
	}

	return &SpacesLowlevelClient{
		token,
		endpoint,
		"v1",
		client,
	}, nil
}

func (c *SpacesLowlevelClient) Get(path string, target interface{}) error {
	url := c.endpoint + "/" + c.version + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Access-Token", c.token)

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 403 {
		return auth.InvalidCredentialsErr{}
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("could not JSON-decode response body: %s", err)
	}

	return nil
}

func (c *SpacesLowlevelClient) Post(path string, body interface{}, target interface{}) error {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := c.endpoint + "/" + c.version + path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Access-Token", c.token)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 403 {
		return auth.InvalidCredentialsErr{}
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("could not JSON-decode response body: %s", err)
	}

	return nil
}