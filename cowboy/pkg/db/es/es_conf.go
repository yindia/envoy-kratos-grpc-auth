package es

import (
	"time"
)

// Conf : holds the params needed to initialize ES client
type Conf struct {
	httpEndpointURLs []string
	rConf            *RetryConf
}

// NewConf : initializes and returns ES Conf struct
func NewConf(httpEndpointURLs []string, rConf *RetryConf) *Conf {
	return &Conf{httpEndpointURLs: httpEndpointURLs, rConf: rConf}
}

// HTTPEndPointURLs : returns the ES http endpoint URLs
func (c *Conf) HTTPEndPointURLs() []string {
	return c.httpEndpointURLs
}

// RConf method : retry config
func (c *Conf) RConf() *RetryConf {
	return c.rConf
}

//RetryConf : holds retry specific config for ES
type RetryConf struct {
	initTimeout time.Duration
	maxTimeout  time.Duration
}

// NewRetryConf : initializes and returns new RetryConf struct
func NewRetryConf(initTimeout time.Duration, maxTimeout time.Duration) *RetryConf {
	return &RetryConf{initTimeout: initTimeout, maxTimeout: maxTimeout}
}

// InitTimeout : method returns the init timeout set. The first retry will
// be made only after this timeout elapses
func (r *RetryConf) InitTimeout() time.Duration {
	return r.initTimeout
}

// MaxTimeout : returns the max timeout set for retries.
// No retries will be made once the max timeout elapses
func (r *RetryConf) MaxTimeout() time.Duration {
	return r.maxTimeout
}
