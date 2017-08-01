package genconfig

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/coldog/k8pack/kubesetup/pkg/signer"
)

type Generator struct {
	// Map from nonce to User.
	users map[string]*User
	lock  sync.Mutex

	Signer *signer.Signer
}

func (gen *Generator) HandleUser(user *User) (string, error) {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	if gen.users == nil {
		gen.users = map[string]*User{}
	}

	nonce := genNonce()
	gen.users[nonce] = user
	return nonce, nil
}

func (gen *Generator) Sign(w http.ResponseWriter, r *http.Request) {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	nonce := r.URL.Query().Get("nonce")

	user, ok := gen.users[nonce]
	if !ok {
		w.WriteHeader(404)
		return
	}

	delete(gen.users, nonce)

	data, err := gen.Signer.SignWithData(signer.CertConfig{
		CN:  user.Name,
		Org: user.Groups,
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
