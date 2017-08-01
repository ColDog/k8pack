package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/genconfig"

	"golang.org/x/oauth2"
)

const (
	DefaultRedirectURL = "http://localhost:6129/"
)

type Providers map[string]*Provider

func (providers Providers) Login(provider string, w http.ResponseWriter, r *http.Request) {
	handler, ok := providers[provider]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("provider not found " + provider))
		return
	}
	handler.Login(w, r)
}

func (providers Providers) Callback(provider string, w http.ResponseWriter, r *http.Request) {
	handler, ok := providers[provider]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("provider not found " + provider))
		return
	}
	handler.Callback(w, r)
}

type FetchUserFunc func(ctx context.Context, token *oauth2.Token) (*genconfig.User, error)

type HandleUserFunc func(user *genconfig.User) (string, error)

type Provider struct {
	Config     *oauth2.Config
	FetchUser  FetchUserFunc
	HandleUser HandleUserFunc
}

func (base *Provider) Login(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"redirect_url": r.URL.Query().Get("redirect_url"),
	}
	state, _ := json.Marshal(data)
	stateStr := hex.EncodeToString(state)

	url := base.Config.AuthCodeURL(stateStr, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (base *Provider) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	data := map[string]string{}

	state := r.FormValue("state")
	if state != "" {
		stateBytes, err := hex.DecodeString(state)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(stateBytes, &data)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
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

	url := data["redirect_url"]
	if url == "" {
		url = DefaultRedirectURL
	}
	url += "?nonce=" + nonce

	http.Redirect(w, r, url, http.StatusFound)
}
