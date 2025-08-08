package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"chaits.org/microservices-repo/pkg/errors"
)

type HTTPClient struct {
	httpclient   *http.Client
	retryOptions *retryOptions
}

type retryOptions struct {
	// MaxRetries is the maximum number of times to retry a failed request.
	MaxRetries int

	// Backoff is the base time duration for exponential backoff between retries.
	Backoff time.Duration

	// RetryableStatusCodes is a map of HTTP status codes that should trigger a retry.
	RetryableStatusCodes map[int]bool
}

type Option func(*HTTPClient)

func NewHTTPClient(opts ...Option) *HTTPClient {
	httpClient := &HTTPClient{
		httpclient:   &http.Client{},
		retryOptions: nil,
	}

	for _, opt := range opts {
		opt(httpClient)
	}
	return httpClient
}

func WithTimeout(timeDuration time.Duration) Option {
	return func(h *HTTPClient) {
		h.httpclient.Timeout = timeDuration
	}
}

func WithRetry(maxRetries int, backOffDuration time.Duration, retryableStatusCodes ...int) Option {
	return func(h *HTTPClient) {
		statusCodeMap := make(map[int]bool)
		for _, statusCode := range retryableStatusCodes {
			statusCodeMap[statusCode] = true
		}
		h.retryOptions = &retryOptions{
			MaxRetries:           maxRetries,
			Backoff:              backOffDuration,
			RetryableStatusCodes: statusCodeMap,
		}
	}
}

func (h *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {

	// No retry options configured, just execute the request once.
	if h.retryOptions == nil {
		return h.httpclient.Do(req)
	}

	var resp *http.Response
	var err error

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	for i := 0; i < h.retryOptions.MaxRetries; i++ {
		if reqBody != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		resp, err = h.httpclient.Do(req)
		if err == nil && !h.retryOptions.RetryableStatusCodes[resp.StatusCode] {
			log.Printf("Not Retrying because Status Code : %d\n", resp.StatusCode)
			return resp, nil
		}

		if i == h.retryOptions.MaxRetries {
			return resp, errors.New(errors.ALL_RETRIES_FAILED_ERROR, fmt.Sprintf("all %d attempts failed, last error: %v", h.retryOptions.MaxRetries, err))
		}

		backOffDuration := h.retryOptions.Backoff * (1 << i)
		log.Printf("Request Failed. (attempt %d of %d). Retrying in %v seconds.", i+1, h.retryOptions.MaxRetries, backOffDuration)
		time.Sleep(backOffDuration)
	}

	return nil, fmt.Errorf("unexpected error after all retries")
}
