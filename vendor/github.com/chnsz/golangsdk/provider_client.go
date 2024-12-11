package golangsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/chnsz/golangsdk/auth"
)

// DefaultUserAgent is the default User-Agent string set in the request header.
const (
	DefaultUserAgent         = "golangsdk/2.0.0"
	DefaultMaxBackoffRetries = 5
)

// UserAgent represents a User-Agent header.
type UserAgent struct {
	// prepend is the slice of User-Agent strings to prepend to DefaultUserAgent.
	// All the strings to prepend are accumulated and prepended in the Join method.
	prepend []string
}

type RetryFunc func(context.Context, *ErrUnexpectedResponseCode, error, uint) error

// Prepend prepends a user-defined string to the default User-Agent string. Users
// may pass in one or more strings to prepend.
func (ua *UserAgent) Prepend(s ...string) {
	ua.prepend = append(s, ua.prepend...)
}

// Join concatenates all the user-defined User-Agend strings with the default
// Gophercloud User-Agent string.
func (ua *UserAgent) Join() string {
	uaSlice := append(ua.prepend, DefaultUserAgent)
	return strings.Join(uaSlice, " ")
}

// ProviderClient stores details that are required to interact with any
// services within a specific provider's API.
//
// Generally, you acquire a ProviderClient by calling the NewClient method in
// the appropriate provider's child package, providing whatever authentication
// credentials are required.
type ProviderClient struct {
	// IdentityBase is the base URL used for a particular provider's identity
	// service - it will be used when issuing authenticatation requests. It
	// should point to the root resource of the identity service, not a specific
	// identity version.
	IdentityBase string

	// IdentityEndpoint is the identity endpoint. This may be a specific version
	// of the identity service. If this is the case, this endpoint is used rather
	// than querying versions first.
	IdentityEndpoint string

	// TokenID is the ID of the most recently issued valid token.
	// NOTE: Aside from within a custom ReauthFunc, this field shouldn't be set by an application.
	// To safely read or write this value, call `Token` or `SetToken`, respectively
	TokenID string

	// ProjectID is the ID of project to which User is authorized.
	ProjectID string

	// DomainID is the ID of project to which User is authorized.
	DomainID string

	// EndpointLocator describes how this provider discovers the endpoints for
	// its constituent services.
	EndpointLocator EndpointLocator

	// HTTPClient allows users to interject arbitrary http, https, or other transit behaviors.
	HTTPClient http.Client

	// UserAgent represents the User-Agent header in the HTTP request.
	UserAgent UserAgent

	// ReauthFunc is the function used to re-authenticate the user if the request
	// fails with a 401 HTTP response code. This a needed because there may be multiple
	// authentication functions for different Identity service versions.
	ReauthFunc func() error

	// AKSKAuthOptions provides the value for AK/SK authentication, it should be nil if you use token authentication,
	// Otherwise, it must have a value
	AKSKAuthOptions AKSKAuthOptions

	// Context is the context passed to the HTTP request.
	Context context.Context

	// Retry backoff func
	RetryBackoffFunc RetryFunc

	// MaxBackoffRetries set the maximum number of backoffs. When not set, defaults to DefaultMaxBackoffRetries
	MaxBackoffRetries uint

	// mut is a mutex for the client. It protects read and write access to client attributes such as getting
	// and setting the TokenID.
	mut *sync.RWMutex

	// reauthmut is a mutex for reauthentication it attempts to ensure that only one reauthentication
	// attempt happens at one time.
	reauthmut *reauthlock
}

// reauthlock represents a set of attributes used to help in the reauthentication process.
type reauthlock struct {
	sync.RWMutex
	reauthing bool
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all
// authenticated service requests.
func (client *ProviderClient) AuthenticatedHeaders() (m map[string]string) {
	if client.reauthmut != nil {
		client.reauthmut.RLock()
		if client.reauthmut.reauthing {
			client.reauthmut.RUnlock()
			return
		}
		client.reauthmut.RUnlock()
	}
	t := client.Token()
	if t == "" {
		return
	}
	return map[string]string{"X-Auth-Token": t}
}

// UseTokenLock creates a mutex that is used to allow safe concurrent access to the auth token.
// If the application's ProviderClient is not used concurrently, this doesn't need to be called.
func (client *ProviderClient) UseTokenLock() {
	client.mut = new(sync.RWMutex)
	client.reauthmut = new(reauthlock)
}

// Token safely reads the value of the auth token from the ProviderClient. Applications should
// call this method to access the token instead of the TokenID field
func (client *ProviderClient) Token() string {
	if client.mut != nil {
		client.mut.RLock()
		defer client.mut.RUnlock()
	}
	return client.TokenID
}

// SetToken safely sets the value of the auth token in the ProviderClient. Applications may
// use this method in a custom ReauthFunc
func (client *ProviderClient) SetToken(t string) {
	if client.mut != nil {
		client.mut.Lock()
		defer client.mut.Unlock()
	}
	client.TokenID = t
}

// Reauthenticate calls client.ReauthFunc in a thread-safe way. If this is
// called because of a 401 response, the caller may pass the previous token. In
// this case, the reauthentication can be skipped if another thread has already
// reauthenticated in the meantime. If no previous token is known, an empty
// string should be passed instead to force unconditional reauthentication.
func (client *ProviderClient) Reauthenticate(previousToken string) (err error) {
	if client.ReauthFunc == nil {
		return nil
	}

	if client.mut == nil {
		return client.ReauthFunc()
	}

	client.mut.Lock()
	defer client.mut.Unlock()

	client.reauthmut.Lock()
	client.reauthmut.reauthing = true
	client.reauthmut.Unlock()

	if previousToken == "" || client.TokenID == previousToken {
		err = client.ReauthFunc()
	}

	client.reauthmut.Lock()
	client.reauthmut.reauthing = false
	client.reauthmut.Unlock()
	return
}

// RequestOpts customizes the behavior of the provider.Request() method.
type RequestOpts struct {
	// JSONBody, if provided, will be encoded as JSON and used as the body of the HTTP request. The
	// content type of the request will default to "application/json" unless overridden by MoreHeaders.
	// It's an error to specify both a JSONBody and a RawBody.
	JSONBody interface{}
	// RawBody contains an io.Reader that will be consumed by the request directly. No content-type
	// will be set unless one is provided explicitly by MoreHeaders.
	RawBody io.Reader
	// JSONResponse, if provided, will be populated with the contents of the response body parsed as
	// JSON.
	JSONResponse interface{}
	// OkCodes contains a list of numeric HTTP status codes that should be interpreted as success. If
	// the response has a different code, an error will be returned.
	OkCodes []int
	// MoreHeaders specifies additional HTTP headers to be provide on the request. If a header is
	// provided with a blank value (""), that header will be *omitted* instead: use this to suppress
	// the default Accept header or an inferred Content-Type, for example.
	MoreHeaders map[string]string
	// ErrorContext specifies the resource error type to return if an error is encountered.
	// This lets resources override default error messages based on the response status code.
	ErrorContext error
	// KeepResponseBody specifies whether to keep the HTTP response body. Usually used, when the HTTP
	// response body is considered for further use. Valid when JSONResponse is nil.
	KeepResponseBody bool
}

// requestState contains temporary state for a single ProviderClient.Request() call.
type requestState struct {
	// This flag indicates if we have reauthenticated during this request because of a 401 response.
	// It ensures that we don't reauthenticate multiple times for a single request. If we
	// reauthenticate, but keep getting 401 responses with the fresh token, reauthenticating some more
	// will just get us into an infinite loop.
	hasReauthenticated bool
	// Retry-After backoff counter, increments during each backoff call
	retries uint
}

var applicationJSON = "application/json"

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)
	return buffer.Bytes(), err
}

// Request performs an HTTP request using the ProviderClient's current HTTPClient. An authentication
// header will automatically be provided.
func (client *ProviderClient) Request(method, url string, options *RequestOpts) (*http.Response, error) {
	return client.doRequest(method, url, options, &requestState{
		hasReauthenticated: false,
	})
}

func (client *ProviderClient) doRequest(method, url string, options *RequestOpts, state *requestState) (*http.Response, error) {
	var body io.Reader
	var contentType *string

	// Derive the content body by either encoding an arbitrary object as JSON, or by taking a provided
	// io.ReadSeeker as-is. Default the content-type to application/json.
	if options.JSONBody != nil {
		if options.RawBody != nil {
			return nil, errors.New("Please provide only one of JSONBody or RawBody to golangsdk.Request()")
		}

		rendered, err := jsonMarshal(options.JSONBody)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(rendered)
		contentType = &applicationJSON
	}

	// Return an error, when "KeepResponseBody" is true and "JSONResponse" is not nil
	if options.KeepResponseBody && options.JSONResponse != nil {
		return nil, errors.New("cannot use KeepResponseBody when JSONResponse is not nil")
	}

	if options.RawBody != nil {
		body = options.RawBody
	}

	// Construct the http.Request.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if client.Context != nil {
		req = req.WithContext(client.Context)
	}

	// Populate the request headers. Apply options.MoreHeaders last, to give the caller the chance to
	// modify or omit any header.
	if contentType != nil {
		req.Header.Set("Content-Type", *contentType)
	}
	req.Header.Set("Accept", applicationJSON)

	// Set the User-Agent header
	req.Header.Set("User-Agent", client.UserAgent.Join())

	if options.MoreHeaders != nil {
		for k, v := range options.MoreHeaders {
			if v != "" {
				req.Header.Set(k, v)
			} else {
				req.Header.Del(k)
			}
		}
	}

	// get latest token from client
	for k, v := range client.AuthenticatedHeaders() {
		req.Header.Set(k, v)
	}

	prereqtok := req.Header.Get("X-Auth-Token")

	if client.AKSKAuthOptions.AccessKey != "" {
		if err := auth.Sign(req, client.AKSKAuthOptions.AccessKey, client.AKSKAuthOptions.SecretKey); err != nil {
			return nil, err
		}
		if client.AKSKAuthOptions.ProjectId != "" {
			req.Header.Set("X-Project-Id", client.AKSKAuthOptions.ProjectId)
		}
		if client.AKSKAuthOptions.DomainID != "" {
			req.Header.Set("X-Domain-Id", client.AKSKAuthOptions.DomainID)
		}
		if client.AKSKAuthOptions.SecurityToken != "" {
			req.Header.Set("X-Security-Token", client.AKSKAuthOptions.SecurityToken)
		}
	}

	// Issue the request.
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Allow default OkCodes if none explicitly set
	okc := options.OkCodes
	if okc == nil {
		okc = defaultOkCodes(method)
	}

	// Validate the HTTP response status.
	var ok bool
	for _, code := range okc {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}

	if !ok {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		respErr := ErrUnexpectedResponseCode{
			URL:       url,
			Method:    method,
			Expected:  okc,
			Actual:    resp.StatusCode,
			RequestId: resp.Header.Get("X-Request-Id"),
			Body:      body,
		}

		errType := options.ErrorContext
		switch resp.StatusCode {
		case http.StatusBadRequest:
			err = ErrDefault400{respErr}
			if error400er, ok := errType.(Err400er); ok {
				err = error400er.Error400(respErr)
			}
		case http.StatusUnauthorized:
			if client.ReauthFunc != nil && !state.hasReauthenticated {
				err = client.Reauthenticate(prereqtok)
				if err != nil {
					e := &ErrUnableToReauthenticate{}
					e.ErrOriginal = respErr
					return nil, e
				}
				if options.RawBody != nil {
					if seeker, ok := options.RawBody.(io.Seeker); ok {
						seeker.Seek(0, 0)
					}
				}
				state.hasReauthenticated = true
				resp, err = client.doRequest(method, url, options, state)
				if err != nil {
					switch err.(type) {
					case *ErrUnexpectedResponseCode:
						e := &ErrErrorAfterReauthentication{}
						e.ErrOriginal = err.(*ErrUnexpectedResponseCode)
						return nil, e
					default:
						e := &ErrErrorAfterReauthentication{}
						e.ErrOriginal = err
						return nil, e
					}
				}
				return resp, nil
			}
			err = ErrDefault401{respErr}
			if error401er, ok := errType.(Err401er); ok {
				err = error401er.Error401(respErr)
			}
		case http.StatusForbidden:
			err = ErrDefault403{respErr}
			if error403er, ok := errType.(Err403er); ok {
				err = error403er.Error403(respErr)
			}
		case http.StatusNotFound:
			err = ErrDefault404{respErr}
			if error404er, ok := errType.(Err404er); ok {
				err = error404er.Error404(respErr)
			}
		case http.StatusMethodNotAllowed:
			err = ErrDefault405{respErr}
			if error405er, ok := errType.(Err405er); ok {
				err = error405er.Error405(respErr)
			}
		case http.StatusRequestTimeout:
			err = ErrDefault408{respErr}
			if error408er, ok := errType.(Err408er); ok {
				err = error408er.Error408(respErr)
			}
		case http.StatusTooManyRequests:
			err = ErrDefault429{respErr}
			if error429er, ok := errType.(Err429er); ok {
				err = error429er.Error429(respErr)
			}

			maxTries := client.MaxBackoffRetries
			if maxTries == 0 {
				maxTries = DefaultMaxBackoffRetries
			}

			if f := client.RetryBackoffFunc; f != nil && state.retries < maxTries {
				var e error

				e = f(client.Context, &respErr, err, state.retries)
				if e != nil {
					return resp, e
				}

				state.retries = state.retries + 1
				return client.doRequest(method, url, options, state)
			}
		case http.StatusInternalServerError:
			err = ErrDefault500{respErr}
			if error500er, ok := errType.(Err500er); ok {
				err = error500er.Error500(respErr)
			}
		case http.StatusBadGateway:
			err = ErrDefault502{respErr}
			if error502er, ok := errType.(Err502er); ok {
				err = error502er.Error502(respErr)
			}
		case http.StatusServiceUnavailable:
			err = ErrDefault503{respErr}
			if error503er, ok := errType.(Err503er); ok {
				err = error503er.Error503(respErr)
			}
		}

		if err == nil {
			err = respErr
		}

		return resp, err
	}

	// Parse the response body as JSON, if requested to do so.
	if options.JSONResponse != nil {
		defer resp.Body.Close()
		// Don't decode JSON when there is no content
		if resp.StatusCode == http.StatusNoContent {
			// read till EOF, otherwise the connection will be closed and cannot be reused
			_, err = io.Copy(ioutil.Discard, resp.Body)
			return resp, err
		}
		if err := json.NewDecoder(resp.Body).Decode(options.JSONResponse); err != nil {
			return nil, err
		}
	}

	// Close unused body to allow the HTTP connection to be reused
	if !options.KeepResponseBody && options.JSONResponse == nil {
		defer resp.Body.Close()
		// read till EOF, otherwise the connection will be closed and cannot be reused
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func defaultOkCodes(method string) []int {
	switch method {
	case "GET", "HEAD":
		return []int{200}
	case "POST":
		return []int{200, 201, 202}
	case "PUT":
		return []int{200, 201, 202}
	case "PATCH":
		return []int{200, 202, 204}
	case "DELETE":
		return []int{200, 202, 204}
	}

	return []int{}
}
