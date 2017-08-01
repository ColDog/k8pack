package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const githubBase = "https://api.github.com"

// NewGithubProvider configures a provider using Github's API. Configuration
// is handled through environment variables.
//
// - client_id     = GITHUB_CLIENT_ID
// - client_secret = GITHUB_CLIENT_SECRET
func NewGithubProvider(advertise string, handleUser HandleUserFunc) *Provider {
	config := &oauth2.Config{
		RedirectURL:  advertise + "/callback/github",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
	return &Provider{
		Config:     config,
		HandleUser: handleUser,
		FetchUser:  fetchGithubUser(config),
	}
}

func fetchGithubUser(config *oauth2.Config) FetchUserFunc {
	return func(ctx context.Context, token *oauth2.Token) (*User, error) {
		resp, err := config.Client(ctx, token).Get(githubBase + "/user")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status)
		}
		data := struct {
			Login string `json:"login"`
			ID    int    `json:"id"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		user := &User{
			Name: data.Login,
		}
		return user, nil
	}
}
