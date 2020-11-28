package es

import (
	"net/http"

	elastic "github.com/olivere/elastic/v7"
)

const maxIdleConns = 100
const maxIdleConnsPerHost = 100

// InitESClient : initializes and returns ES client
func InitESClient(conf *Conf) (*elastic.Client, error) {
	httpClient := getDefaultHTTPClient()
	rConf := conf.RConf()

	var r elastic.Retrier
	if rConf != nil {
		b := elastic.NewExponentialBackoff(rConf.InitTimeout(), rConf.MaxTimeout())
		r = elastic.NewBackoffRetrier(b)
	}

	client, err := elastic.NewClient(elastic.SetHttpClient(httpClient), elastic.SetURL(conf.HTTPEndPointURLs()...), elastic.SetSniff(false), elastic.SetRetrier(r))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getDefaultHTTPClient : return http client with default params
func getDefaultHTTPClient() *http.Client {
	defaultTransportPointer, _ := http.DefaultTransport.(*http.Transport)
	defaultTransport := defaultTransportPointer
	defaultTransport.MaxIdleConns = maxIdleConns
	defaultTransport.MaxIdleConnsPerHost = maxIdleConnsPerHost
	defaultHTTPClient := &http.Client{Transport: defaultTransport}

	return defaultHTTPClient
}
