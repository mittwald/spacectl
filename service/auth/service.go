package auth

import (
	"errors"
	"encoding/json"
	"net/http"
	"bytes"
	"time"
	"fmt"
)

type authenticationRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

type AuthenticationResult struct {
	Token string
	ValidUntil time.Time
}

type AuthenticationService struct {
	AuthServerURL string
}

func (a *AuthenticationService) Authenticate(emailAddress string, password string) (*AuthenticationResult, error) {
	if emailAddress == "" {
		return nil, errors.New("empty username")
	}

	if password == "" {
		return nil, errors.New("empty password")
	}

	body, err := json.Marshal(authenticationRequest{
		EmailAddress: emailAddress,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", a.AuthServerURL + "/v1/auth/local/login", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 403 {
		return nil, InvalidCredentialsErr{}
	}

	//resBody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, err
	//}

	if res.StatusCode >= 400 {
		return nil, AuthErr{fmt.Errorf("unexpected status code: %d", res.StatusCode)}
	}

	authResponse := AuthenticationResult{}
	err = json.NewDecoder(res.Body).Decode(&authResponse)
	if err != nil {
		return nil, AuthErr{fmt.Errorf("could not JSON-decode response body: %s", err)}
	}

	return &authResponse, nil
}