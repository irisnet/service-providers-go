package monitor

import "time"

type IMonitor interface {
	Start()
	Stop()
}

type Endpoint struct {
	URL  string
	Auth Authentication
}

func NewEndpoint(url, accessKey, secret string) Endpoint {
	auth := Authentication{
		AccessKey: accessKey,
		Secret:    secret,
	}

	return Endpoint{
		URL:  url,
		Auth: auth,
	}
}

func NewEndpointFromURL(url string) Endpoint {
	return Endpoint{
		URL: url,
	}
}

type Authentication struct {
	AccessKey string
	Secret    string
}

type RetryConfig struct {
	Timeout  time.Duration
	Attempts int
}

func NewRetryConfig(timeout time.Duration, attempts int) RetryConfig {
	return RetryConfig{
		Timeout:  timeout,
		Attempts: attempts,
	}
}
