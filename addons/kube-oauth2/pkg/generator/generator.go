package generator

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/provider"
	"github.com/coldog/k8pack/kubesetup/pkg/signer"
)

const expireAfter = 5 * time.Minute

type item struct {
	user    *provider.User
	expires time.Time
}

// Generator generates Kubeconfig files for users.
type Generator struct {
	// Map from nonce to User.
	store map[string]*item
	lock  sync.Mutex

	Signer *signer.Signer
}

// HandleUser stores the user and returns a nonce to retrieve this user.
func (gen *Generator) HandleUser(user *provider.User) (string, error) {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	if gen.store == nil {
		gen.store = map[string]*item{}
	}

	nonce := genNonce()
	gen.store[nonce] = &item{
		expires: time.Now().Add(expireAfter),
		user:    user,
	}
	return nonce, nil
}

// Sign is an http handler to create a signed Kubeconfig.
func (gen *Generator) Sign(w http.ResponseWriter, r *http.Request) {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	nonce := r.URL.Query().Get("nonce")

	item, ok := gen.store[nonce]
	if !ok {
		w.WriteHeader(404)
		return
	}
	delete(gen.store, nonce)

	if time.Now().After(item.expires) {
		w.WriteHeader(404)
		return
	}

	data, err := gen.Signer.SignWithData(signer.CertConfig{
		CN:  item.user.Name,
		Org: item.user.Groups,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}

func genNonce() string {
	return strconv.FormatInt(rand.Int63(), 10)
}
