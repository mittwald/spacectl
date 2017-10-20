package lowlevel

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/mittwald/spacectl/service/auth"
	"bytes"
	"time"
	"log"
	"regexp"
	"strings"
	"io"
	"io/ioutil"
)

var versionRegexp = regexp.MustCompile("^/v[0-9]+/")

type SpacesLowlevelClient struct {
	token string
	endpoint string
	version string

	client *http.Client
	logger *log.Logger
}

func NewSpacesLowlevelClient(token string, endpoint string, logger *log.Logger) (*SpacesLowlevelClient, error) {
	client := &http.Client{
	}

	return &SpacesLowlevelClient{
		token,
		endpoint,
		"v1",
		client,
		logger,
	}, nil
}

func (c *SpacesLowlevelClient) pathWithVersion(path string) string {
	if versionRegexp.MatchString(path) {
		return path
	}

	return "/" + c.version + "/" + strings.TrimPrefix(path, "/")
}

func (c *SpacesLowlevelClient) Get(path string, target interface{}) error {
	url := c.endpoint + c.pathWithVersion(path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Access-Token", c.token)

	c.logger.Printf("executing GET on %s", url)

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
		msg := Message{}

		// The error here can safely be ignored since it does not matter much, anyway.
		// Either the response body contains a "msg" or it doesn't.
		_ = json.NewDecoder(res.Body).Decode(&msg)

		return ErrUnexpectedStatusCode{res.StatusCode, msg.String()}
	}

	var buf bytes.Buffer
	io.Copy(&buf, res.Body)

	reader := bytes.NewReader(buf.Bytes())

	err = json.NewDecoder(reader).Decode(target)
	if err != nil {
		return fmt.Errorf("could not JSON-decode response body: %s", err)
	}

	reader.Seek(0, io.SeekStart)
	responseBytes, _ := ioutil.ReadAll(reader)
	c.logger.Println(string(responseBytes))

	return nil
}

func (c *SpacesLowlevelClient) Post(path string, body interface{}, target interface{}) error {
	return c.request("POST", path, body, target)
}

func (c *SpacesLowlevelClient) Put(path string, body interface{}, target interface{}) error {
	return c.request("PUT", path, body, target)
}

func (c *SpacesLowlevelClient) Delete(path string, target interface{}) error {
	return c.request("DELETE", path, nil, target)
}

func (c *SpacesLowlevelClient) GetCanonicalURL(path string) (string, error) {
	url := c.endpoint + c.pathWithVersion(path)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Access-Token", c.token)

	c.logger.Printf("executing HEAD on %s", url)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 300 {
		return req.URL.String(), nil
	} else if res.StatusCode >= 300 && res.StatusCode < 400 {
		url, err := res.Location()
		if err != nil {
			return "", err
		}
		return url.String(), nil
	}

	return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
}

func (c *SpacesLowlevelClient) request(method string, path string, body interface{}, target interface{}) error {
	var reqBody []byte = []byte{}
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	url := c.endpoint + c.pathWithVersion(path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	req.Header.Set("X-Access-Token", c.token)

	c.logger.Printf("executing %s on %s: %s", method, url, string(reqBody))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	c.logger.Printf("response code: %d", res.StatusCode)

	if res.StatusCode == 403 {
		return auth.InvalidCredentialsErr{}
	}

	if res.StatusCode >= 400 {
		msg := Message{}

		var buf bytes.Buffer
		io.Copy(&buf, res.Body)

		reader := bytes.NewReader(buf.Bytes())

		// The error here can safely be ignored since it does not matter much, anyway.
		// Either the response body contains a "msg" or it doesn't.
		//_ = json.NewDecoder(res.Body).Decode(&msg)
		json.Unmarshal(buf.Bytes(), &msg)

		reader.Seek(0, io.SeekStart)
		responseBytes, _ := ioutil.ReadAll(reader)
		c.logger.Println(string(responseBytes))

		c.logger.Printf("response: %v", msg)

		return ErrUnexpectedStatusCode{res.StatusCode, msg.String()}
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("could not JSON-decode response body: %s", err)
	}

	c.logger.Printf("response: %s", target)

	return nil
}