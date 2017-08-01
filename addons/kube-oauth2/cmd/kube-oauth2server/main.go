package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/coldog/k8pack/kubesetup/pkg/signer"

	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/genconfig"
	"github.com/coldog/k8pack/addons/kube-oauth2/pkg/provider"

	"github.com/julienschmidt/httprouter"
)

var (
	advertise string
	listen    string
	caURI     string
	keyURI    string
	api       string
	cluster   string

	githubEnabled bool
)

type server struct {
	providers provider.Providers
	generator *genconfig.Generator
}

func (s *server) login(w http.ResponseWriter, r *http.Request, rp httprouter.Params) {
	s.providers.Login(rp.ByName("provider"), w, r)
}

func (s *server) callback(w http.ResponseWriter, r *http.Request, rp httprouter.Params) {
	s.providers.Callback(rp.ByName("provider"), w, r)
}

func (s *server) sign(w http.ResponseWriter, r *http.Request, rp httprouter.Params) {
	s.generator.Sign(w, r)
}

func main() {
	flag.StringVar(&advertise, "advertise", "http://127.0.0.1:3000", "Host address")
	flag.StringVar(&listen, "listen", "127.0.0.1:3000", "Listen address")
	flag.StringVar(&caURI, "ca-uri", "", "CA certificate uri")
	flag.StringVar(&keyURI, "key-uri", "", "Secret key uri")
	flag.StringVar(&api, "api", "", "Kubernetes api host")
	flag.StringVar(&cluster, "cluster", "default", "Cluster name")
	flag.BoolVar(&githubEnabled, "github", false, "Enable github auth plugin")
	flag.Parse()

	sign := &signer.Signer{
		CaURI:   caURI,
		KeyURI:  keyURI,
		APIHost: api,
		Cluster: cluster,
	}
	err := sign.Init()
	if err != nil {
		fmt.Println("failed to init signer", err)
		os.Exit(1)
	}

	srv := &server{
		providers: provider.Providers{},
		generator: &genconfig.Generator{
			Signer: sign,
		},
	}

	if githubEnabled {
		srv.providers["github"] = provider.NewGithubProvider(advertise, srv.generator.HandleUser)
	}

	r := httprouter.New()

	r.GET("/sign", srv.sign)
	r.GET("/login/:provider", srv.login)
	r.GET("/callback/:provider", srv.callback)

	handler := requestLogger{r}

	fmt.Println("serving on", listen)
	if err := http.ListenAndServe(listen, handler); err != nil {
		fmt.Println("failed serve", err)
		os.Exit(1)
	}
}
