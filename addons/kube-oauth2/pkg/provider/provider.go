package provider

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/genconfig"

	"golang.org/x/oauth2"
)

var (
	Secret []byte = []byte("secret")
)

const (
	DefaultRedirectURL = "http://localhost:6129/"
)

type User struct {
	Name   string
	Groups []string
}

// Providers allows for multiple named providers. It calls the login and callback
// methods on each named provider.
type Providers map[string]*Provider

// Login redirects to the named oauth2 provider.
func (providers Providers) Login(provider string, w http.ResponseWriter, r *http.Request) {
	handler, ok := providers[provider]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("provider not found " + provider))
		return
	}
	handler.Login(w, r)
}

// Callback handles the callback from the named oauth2 provider.
func (providers Providers) Callback(provider string, w http.ResponseWriter, r *http.Request) {
	handler, ok := providers[provider]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("provider not found " + provider))
		return
	}
	handler.Callback(w, r)
}

// FetchUserFunc should fetch a user
type FetchUserFunc func(ctx context.Context, token *oauth2.Token) (*genconfig.User, error)

type HandleUserFunc func(user *genconfig.User) (string, error)

// Provider is an implementation of the two core oauth2 routes, login and the callback.
// The callback route will call the fetch user handler and then the handle user handler
// these should return and then configure the user appropriately.
type Provider struct {
	Config     *oauth2.Config
	FetchUser  FetchUserFunc
	HandleUser HandleUserFunc
}

// Login redirects to the oauth2 provider.
func (base *Provider) Login(w http.ResponseWriter, r *http.Request) {
	state := encodeState(map[string]string{
		"redirect_url": r.URL.Query().Get("redirect_url"),
	})
	url := base.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// Callback handles the callback from the oauth2 provider.
func (base *Provider) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	state, err := decodeState(r.FormValue("state"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	errorMsg := r.FormValue("error")
	if errorMsg != "" {
		w.WriteHeader(500)
		w.Write([]byte(errorMsg))
		return
	}

	code := r.FormValue("code")
	token, err := base.Config.Exchange(ctx, code)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := base.FetchUser(r.Context(), token)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	nonce, err := base.HandleUser(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	url := state["redirect_url"]
	if url == "" {
		url = DefaultRedirectURL
	}
	url += "?nonce=" + nonce

	http.Redirect(w, r, url, http.StatusFound)
}

func encodeState(m map[string]string) string {
	state, _ := json.Marshal(m)
	encoded := hex.EncodeToString(state)
	return encoded + "." + computeSignature(encoded)
}

func decodeState(state string) (map[string]string, error) {
	if state == "" {
		return nil, errors.New("invalid state")
	}

	spl := strings.Split(state, ".")
	if len(spl) < 2 {
		return nil, errors.New("invalid state")
	}

	state = spl[0]
	signature := spl[1]

	if signature != computeSignature(state) {
		return nil, errors.New("invalid state signature")
	}

	stateBytes, err := hex.DecodeString(state)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %v", err)
	}

	data := map[string]string{}
	err = json.Unmarshal(stateBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %v", err)
	}
	return data, nil
}

func computeSignature(data string) string {
	m := hmac.New(sha1.New, Secret)
	return hex.EncodeToString(m.Sum([]byte(data)))
}
