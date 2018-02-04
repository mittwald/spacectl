package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/mittwald/spacectl/system/browser"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
)

type OAuthAuthenticationService struct {
	AuthServerURL string
}

func (a *OAuthAuthenticationService) Authenticate() (*AuthenticationResult, error) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID: "spaces.de/oauth/spacectl",
		ClientSecret: "",
		Scopes: []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL: fmt.Sprintf("%s/o/oauth2/auth", a.AuthServerURL),
			TokenURL: fmt.Sprintf("%s/o/oauth2/token", a.AuthServerURL),
		},
		RedirectURL: "http://localhost:6241/oauth-redir",
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	url := conf.AuthCodeURL(token)

	accessTokenChan := make(chan *AuthenticationResult)
	errChan := make(chan error)

	t := template.Must(template.New("error").Parse(errorPage))

	go func() {
		http.HandleFunc("/oauth-redir", func(rw http.ResponseWriter, req *http.Request) {
			fmt.Println("HTTP callback called")

			query := req.URL.Query()
			authorizationCode := query.Get("code")

			token, err := conf.Exchange(ctx, authorizationCode)
			if err != nil {
				errChan <- err
				rw.WriteHeader(500)
				t.Execute(rw, map[string]string{"error": err.Error()})
				return
			}

			result := &AuthenticationResult{
				Token: token.AccessToken,
				ValidUntil: token.Expiry,
			}

			accessTokenChan <- result

			rw.Write([]byte(confirmationPage))
		})
		http.ListenAndServe("127.0.0.1:6241", nil)
	}()

	fmt.Println("Opening the OAuth authorization page in your browser.")
	fmt.Println("If your browser does not open, visit the following URL:")
	fmt.Println("  " + url)

	browser.OpenURLFork(url)

	fmt.Println("waiting for HTTP callback")

	select {
	case r := <-accessTokenChan:
		return r, nil
	case err := <-errChan:
		return nil, err
	}
}
