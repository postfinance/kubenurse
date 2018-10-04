package checker

import (
	"errors"
	"net/http"
)

// doRequest does an http request only to get the http status code
func (c *Checker) doRequest(url string) (string, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return err.Error(), err
	}

	// Body is non-nil if err is nil, so close it
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "ok", nil
	}

	return resp.Status, errors.New(resp.Status)
}
