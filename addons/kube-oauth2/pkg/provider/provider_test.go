package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
)

var tokenURL string

func init() {
	sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			RefreshToken string `json:"refresh_token"`
		}{AccessToken: "test"})
	}))
	tokenURL = sv.URL
}

func TestProvider_Login(t *testing.T) {
	p := &Provider{
		Config: &oauth2.Config{},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login", nil)

	p.Login(w, r)

	assert.Equal(t, http.StatusFound, w.Code)
}

func TestProvider_Callback(t *testing.T) {
	handleUser := func(user *genconfig.User) (string, error) {
		return "", nil
	}
	fetchUser := func(ctx context.Context, token *oauth2.Token) (*genconfig.User, error) {
		return &genconfig.User{}, nil
	}

	p := &Provider{
		HandleUser: handleUser,
		FetchUser:  fetchUser,
		Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{TokenURL: tokenURL},
		},
	}

	state := encodeState(map[string]string{"test": ""})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/callback?state="+state, nil)

	p.Callback(w, r)

	assert.Equal(t, http.StatusFound, w.Code)
}

func TestProvider_CallbackNoState(t *testing.T) {
	handleUser := func(user *genconfig.User) (string, error) {
		return "", nil
	}
	fetchUser := func(ctx context.Context, token *oauth2.Token) (*genconfig.User, error) {
		return &genconfig.User{}, nil
	}

	p := &Provider{
		HandleUser: handleUser,
		FetchUser:  fetchUser,
		Config: &oauth2.Config{
			Endpoint: oauth2.Endpoint{TokenURL: tokenURL},
		},
	}

	state := "bleh"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/callback?state="+state, nil)

	p.Callback(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
