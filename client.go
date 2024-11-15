package gobalena

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	APIRetryCount      = 3
	APIRetryBackoff    = 5 * time.Second
	MaxAPIRetryBackoff = 30 * time.Second
)

type SturdyClient struct {
	*resty.Client
}

func NewSturdyHTTPClient() *SturdyClient {
	return &SturdyClient{
		resty.New().
			SetHeader("Content-Type", "application/json").
			SetRetryCount(APIRetryCount).
			SetRetryWaitTime(APIRetryBackoff).
			SetDebug(true).
			EnableGenerateCurlOnDebug().
			EnableTrace().
			SetRetryMaxWaitTime(MaxAPIRetryBackoff).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				if err != nil {
					// If there's an error at the network level, always retry
					return true
				}

				statusCode := r.StatusCode()
				// Retry for 5XX server errors or 429 Too Many Requests
				if statusCode == 429 || (statusCode >= 500 && statusCode <= 599) {
					return true
				}

				// Do not retry on 4XX errors other than 429
				return false
			}),
	}
}

func (c *SturdyClient) SetRetryCount(count int) *SturdyClient {
	c.Client.SetRetryCount(count)
	return c
}

func (c *SturdyClient) SetRetryWaitTime(waitTime time.Duration) *SturdyClient {
	c.Client.SetRetryWaitTime(waitTime)
	return c
}

func (c *SturdyClient) SetRetryMaxWaitTime(maxWaitTime time.Duration) *SturdyClient {
	c.Client.SetRetryMaxWaitTime(maxWaitTime)
	return c
}

func (c *SturdyClient) SetHeader(header, value string) *SturdyClient {
	c.Client.SetHeader(header, value)
	return c
}

func (c *SturdyClient) SetBaseURL(url string) *SturdyClient {
	c.Client.SetBaseURL(url)
	return c
}
