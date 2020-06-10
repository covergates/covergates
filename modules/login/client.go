package login

import (
	"crypto/tls"
	"net/http"
)

func BasicClient(insecure bool) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}
	return client
}
