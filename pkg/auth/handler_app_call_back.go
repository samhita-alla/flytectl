package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

func CallbackHandler(c oauth2.Config) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		codeVerifier := resetPKCE()
		if req.URL.Query().Get("error") != "" {
			TokenChannel <- fmt.Sprintf("error on callback during authorization due to %v", req.URL.Query().Get("error"))
			return
		}
		if req.URL.Query().Get("code") == "" {
			TokenChannel <- fmt.Sprint("Could not find the authorize code")
			return
		}
		// We'll check whether we sent a code+PKCE request, and if so, send the code_verifier along when requesting the access token.
		var opts []oauth2.AuthCodeOption
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))

		token, err := c.Exchange(context.Background(), req.URL.Query().Get("code"), opts...)
		if err != nil {
			TokenChannel <- fmt.Sprintf("error while exchanging auth code due to %v", err)
			return
		}
		TokenChannel <- token.AccessToken
	}
}

