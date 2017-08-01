package generator

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/provider"
	"github.com/coldog/k8pack/kubesetup/pkg/signer"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_HandleUser(t *testing.T) {
	g := &Generator{}

	nonce, err := g.HandleUser(&provider.User{
		Name: "test",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, nonce)
}

func TestGenerator_Sign(t *testing.T) {
	g := &Generator{
		Signer: &signer.Signer{
			CaURI:   "file://./test-ca.pem",
			KeyURI:  "file://./test-ca-key.pem",
			APIHost: "k8s.test.ca",
			Cluster: "default",
		},
	}

	err := g.Signer.Init()
	assert.Nil(t, err)

	nonce, err := g.HandleUser(&provider.User{
		Name: "test",
	})
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/sign?nonce="+nonce, nil)

	g.Sign(w, r)

	assert.Equal(t, 200, w.Code)

	io.Copy(os.Stdout, w.Body)
}
