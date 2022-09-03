package readers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	caCertFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Reader struct {
	token      string
	client     *http.Client
	partialUrl string
}

func (r *Reader) GetStatus(namespace string) (string, int) {
	req, err := http.NewRequest("GET", r.partialUrl+namespace+"/pods", nil)
	if err != nil {
		log.Printf("Failed to create new request: %v\n", err)
		return "Unexpected Error", 500
	}

	req.Header.Set("Authorization", "Bearer "+r.token)
	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("Request failed with error: %v\n", err)
		return "Unexpected Error", 500
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), 200
}

func NewReader() *Reader {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Fatalf("Error getting serviceaccount token: %v\n", err)
	}

	var kubeApiHost, kubeApiPort string
	var exists bool
	kubeApiHost, exists = os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !exists {
		log.Println("Could not find KUBERNETES_SERVICE_HOST. Defaulting to localhost")
		kubeApiHost = "localhost"
	}
	kubeApiPort, exists = os.LookupEnv("KUBERNETES_SERVICE_PORT")
	if !exists {
		log.Println("Could not find KUBERNETES_SERVICE_PORT. Defaulting to 6443")
		kubeApiPort = "6443"
	}

	partialUrl := fmt.Sprintf("https://%s:%s/api/v1/namespaces/", kubeApiHost, kubeApiPort)

	var caCert []byte
	caCert, err = ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error getting kube ca cert: %v\n", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	return &Reader{
		token:      strings.Trim(string(token), "\n"),
		client:     &client,
		partialUrl: partialUrl,
	}
}
