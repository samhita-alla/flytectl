package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	goauth "golang.org/x/oauth2"
	"google.golang.org/grpc"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// The following provides the setup required for the client to perform the "Authorization Code" flow with PKCE in order
// to obtain an access token for public/untrusted clients.

const cookiePKCE = "isPKCE"

var (
	// pkceCodeVerifier stores the generated random value which the client will on-send to the auth server with the received
	// authorization code. This way the oauth server can verify that the base64URLEncoded(sha265(codeVerifier)) matches
	// the stored code challenge, which was initially sent through with the code+PKCE authorization request to ensure
	// that this is the original user-agent who requested the access token.
	PkceCodeVerifier string

	// pkceCodeChallenge stores the base64(sha256(codeVerifier)) which is sent from the
	// client to the auth server as required for PKCE.
	PkceCodeChallenge string


	TokenChannel chan string
)

// The following sets up the requirements for generating a standards compliant PKCE code verifier.
const codeVerifierLenMin = 43
const codeVerifierLenMax = 128
const codeVerifierAllowedLetters = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ._~"

// generateCodeVerifier provides an easy way to generate an n-length randomised
// code verifier.
func generateCodeVerifier(n int) string {
	// Enforce standards compliance...
	if n < codeVerifierLenMin {
		n = codeVerifierLenMin
	}
	if n > codeVerifierLenMax {
		n = codeVerifierLenMax
	}

	// Randomly choose some allowed characters...
	b := make([]byte, n)
	for i := range b {
		// ensure we use non-deterministic random ints.
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeVerifierAllowedLetters))))
		b[i] = codeVerifierAllowedLetters[j.Int64()]
	}

	return string(b)
}

// generateCodeChallenge returns a standards compliant PKCE S(HA)256 code
// challenge.
func generateCodeChallenge(codeVerifier string) string {
	// Create a sha-265 hash from the code verifier...
	s256 := sha256.New()
	s256.Write([]byte(codeVerifier))

	// Then base64 encode the hash sum to create a code challenge...
	return base64.RawURLEncoding.EncodeToString(s256.Sum(nil))
}

// resetPKCE cleans up PKCE details and returns the code verifier.
func resetPKCE() (codeVerifier string) {
	codeVerifier = PkceCodeVerifier
	PkceCodeVerifier = ""
	return codeVerifier
}


// A valid oauth2 client (check the store) that additionally requests an OpenID Connect id token
var clientConf = goauth.Config{
	ClientID:     "flytectl",
	ClientSecret: "foobar",
	RedirectURL:  "http://localhost:3846/callback",
	Scopes:       []string{"all"},
	Endpoint: goauth.Endpoint{
		TokenURL: "http://localhost:8088/oauth2/token",
		AuthURL:  "http://localhost:8088/oauth2/authorize",
	},
}

type FlyteCtlTokenSource struct {
	accessToken string
}

func (ts *FlyteCtlTokenSource) Token() (*goauth.Token, error) {
	t := &goauth.Token{
		AccessToken: ts.accessToken,
		Expiry:      time.Now().Add(1 * time.Minute),
		TokenType:   "bearer",
	}
	return t, nil
}

type InsecurePerRPCCredentials struct{
	goauth.TokenSource
}

func (cr InsecurePerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := cr.Token()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (cr InsecurePerRPCCredentials) RequireTransportSecurity() bool {
	return false
}

func StartAuthFlow(ctx context.Context) (grpc.CallOption, error) {
	// ### oauth2 client ###
	http.HandleFunc("/callback", CallbackHandler(clientConf)) // the oauth2 callback endpoint
	port := "3846"
	TokenChannel = make(chan string)
	if os.Getenv("FLYTE_AUTH_PORT") != "" {
		port = os.Getenv("FLYTE_AUTH_PORT")
	}
	PkceCodeVerifier = generateCodeVerifier(64)
	PkceCodeChallenge = generateCodeChallenge(PkceCodeVerifier)
	urlToOpen := clientConf.AuthCodeURL("some-random-state-foobar") + "&nonce=some-random-nonce&code_challenge=" +
		PkceCodeChallenge + "&code_challenge_method=S256"
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			logger.Fatal(ctx, "Couldn't start the callback http server on port %v", port)
		}
	}()
	fmt.Println("Please open your webbrowser at " + urlToOpen)
	_ = exec.Command("open", urlToOpen).Run()
	token := <- TokenChannel
	var callOption grpc.CallOption
	accessToken := FlyteCtlTokenSource{
		accessToken: token,
	}
	callOption = grpc.PerRPCCredsCallOption{Creds: InsecurePerRPCCredentials{TokenSource: &accessToken}}
	return callOption, nil
}

type GetExecution func(context.Context, *admin.WorkflowExecutionGetRequest, ...grpc.CallOption) (*admin.Execution, error)

func OauthGetExecutionCallDecorator(getExecGrpcCall GetExecution) GetExecution {
	return func(ctx context.Context, msg *admin.WorkflowExecutionGetRequest, callOptions ...grpc.CallOption) (*admin.Execution, error) {
		result, err := getExecGrpcCall(ctx, msg, callOptions ...)
		if err != nil {
			var authFlowCallOption grpc.CallOption
			if authFlowCallOption, err = StartAuthFlow(ctx); err != nil {
				return nil, err
			}
			callOptions = append(callOptions, authFlowCallOption)
			result, err = getExecGrpcCall(ctx, msg, callOptions ...)
		}
		return result, err
	}
}


type ListExecution func(context.Context, *admin.ResourceListRequest, ...grpc.CallOption) (*admin.ExecutionList, error)

func OauthListExecutionCallDecorator(listExecGrpcCall ListExecution) ListExecution {
	return func(ctx context.Context, msg *admin.ResourceListRequest, callOptions ...grpc.CallOption) (*admin.ExecutionList, error) {
		result, err := listExecGrpcCall(ctx, msg, callOptions ...)
		if err != nil {
			var authFlowCallOption grpc.CallOption
			if authFlowCallOption, err = StartAuthFlow(ctx); err != nil {
				return nil, err
			}
			callOptions = append(callOptions, authFlowCallOption)
			result, err = listExecGrpcCall(ctx, msg, callOptions ...)
		}
		return result, err
	}
}

