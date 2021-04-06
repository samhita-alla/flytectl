package auth

import (
	"context"
	"fmt"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/pkg/browser"
	goauth "golang.org/x/oauth2"
	"google.golang.org/grpc"
	"net/http"
	"net/url"
	"time"
)

const (
	Timeout = 15 * time.Second
)

func StartAuthFlow(ctx context.Context) (grpc.CallOption, error) {
	var clientConf goauth.Config
	var err error
	// Generate the client config by fetching the discovery endpoint data from admin.
	if clientConf, err = GenerateClientConfig(ctx); err != nil {
		return nil, err
	}
	var redirectUrl *url.URL
	if redirectUrl, err = url.Parse(clientConf.RedirectURL); err != nil {
		return nil, err
	}
	// Register the call back handler
	http.HandleFunc(redirectUrl.Path, callbackHandler(clientConf)) // the oauth2 callback endpoint

	tokenChannel = make(chan *goauth.Token, 1)
	errorChannel = make(chan error, 1)
	timeoutChannel = make(chan bool, 1)
	// Run timeout go routine inorder to timeout the authflow incase there are no redirects on the http endpoint created by the app
	go func() {
		time.Sleep(Timeout)
		timeoutChannel <- true
	}()

	pkceCodeVerifier = generateCodeVerifier(64)
	pkceCodeChallenge = generateCodeChallenge(pkceCodeVerifier)
	stateString = state(32)
	nonces = state(32)
	// Replace S256 with one from cient config and provide a support to generate code challenge using the passed in method.
	urlToOpen := clientConf.AuthCodeURL(stateString) + "&nonce=" + nonces + "&code_challenge=" +
		pkceCodeChallenge + "&code_challenge_method=S256"

	go func() {
		if err = http.ListenAndServe(redirectUrl.Host, nil); err != nil {
			logger.Fatal(ctx, "Couldn't start the callback http server on host %v due to %v", redirectUrl.Host, err)
		}
	}()
	fmt.Println("Opening the browser at " + urlToOpen)
	if err = browser.OpenURL(urlToOpen); err != nil {
		return nil, err
	}
	var token *goauth.Token
	select {
	case err = <-errorChannel:
		return nil, err
	case _ = <-timeoutChannel:
		return nil, fmt.Errorf("timeout occured during auth flow")
	case token = <-tokenChannel:
		var callOption grpc.CallOption
		accessToken := FlyteCtlTokenSource{
			flyteCtlToken: token,
		}
		callOption = grpc.PerRPCCredsCallOption{Creds: InsecurePerRPCCredentials{TokenSource: &accessToken}}
		return callOption, nil
	}
}
