package auth

import (
	"errors"
	"k8s.io/client-go/pkg/util/json"
	"net/http"
	"bytes"
	"time"
	"io/ioutil"
	"fmt"
)

type authenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationResult struct {
	Token string
	ValidUntil time.Time
}

type AuthenticationService struct {
	AuthServerURL string
}

func (a *AuthenticationService) Authenticate(username string, password string) (string, error) {
	if username == "" {
		return "", errors.New("empty username")
	}

	if password == "" {
		return "", errors.New("empty password")
	}

	body, err := json.Marshal(authenticationRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", a.AuthServerURL + "/v1/authenticate", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == 403 {
		return "", InvalidCredentialsErr{}
	}
	if res.StatusCode >= 400 {
		return "", AuthErr{fmt.Errorf("unexpected status code: %d", res.StatusCode)}
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}