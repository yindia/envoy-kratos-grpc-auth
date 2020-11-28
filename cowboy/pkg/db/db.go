package db

import (
	"os"
	"strings"
	"time"

	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/db/es"

	elastic "github.com/olivere/elastic/v7"
)

// ToDo : use a factory when more datastores included

// NewESClient : returns instance of es client
func NewESClient() (*elastic.Client, error) {
	esEndPoint := strings.Split(os.Getenv("ELS_URL"), ",")
	rConf := es.NewRetryConf(2*time.Second, 10*time.Second)
	eConf := es.NewConf(esEndPoint, rConf)
	return es.InitESClient(eConf)
}
