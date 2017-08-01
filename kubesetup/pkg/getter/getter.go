package getter

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Get fetches a file or http resource depending on the prefix.
func Get(uri string) ([]byte, error) {
	if strings.HasPrefix(uri, "file://") {
		uri = strings.Replace(uri, "file://", "", 1)
		return ioutil.ReadFile(uri)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
