package httpclient_go

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

type HttpClientGo struct {
	signer    *Signer
	Method    string
	Url       string
	Body      interface{}
	Header    map[string]string
	request   *http.Request
	Error     error
	Transport *http.Transport
}

func NewHttpClientGo(c *config.Config) (*HttpClientGo, diag.Diagnostics) {

	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, diag.Errorf("AKSK is not set")
	}

	return &HttpClientGo{
		signer: &Signer{
			Key:    c.AccessKey,
			Secret: c.SecretKey,
		},
		Header: map[string]string{
			"content-type": "application/json",
		},
	}, nil
}

func (client *HttpClientGo) WithMethod(method string) *HttpClientGo {
	client.Method = method
	return client
}

func (client *HttpClientGo) WithUrl(url string) *HttpClientGo {
	client.Url = url
	return client
}

func (client *HttpClientGo) WithUrlWithoutEndpoint(cfg *config.Config, srv, region, path string) *HttpClientGo {
	endpoint := config.GetServiceEndpoint(cfg, srv, region)

	client.Url = endpoint + path
	return client
}

func (client *HttpClientGo) WithBody(body interface{}) *HttpClientGo {
	client.Body = body
	return client
}

func (client *HttpClientGo) WithHeader(header map[string]string) *HttpClientGo {
	if len(header) == 0 {
		return client
	}
	client.Header = header
	return client
}

func (client *HttpClientGo) ToRequest() {
	var err error
	if client.Body == nil {
		client.request, err = http.NewRequest(client.Method, client.Url, nil)
		if err != nil {
			client.Error = err
			return
		}
	} else {
		b, err := json.Marshal(client.Body)
		if err != nil {
			client.Error = err
			return
		}
		client.request, err = http.NewRequest(client.Method, client.Url, ioutil.NopCloser(bytes.NewBuffer(b)))
		if err != nil {
			client.Error = err
			return
		}
	}
	for k, v := range client.Header {
		client.request.Header.Add(k, v)
	}
}

func (client *HttpClientGo) Do() (*http.Response, error) {
	if client.request == nil {
		client.ToRequest()
	}
	err := client.signer.Sign(client.request)
	if err != nil {
		return nil, err
	}
	c := http.DefaultClient
	if client.Transport != nil {
		c.Transport = client.Transport
	}
	return c.Do(client.request)
}

func (client *HttpClientGo) WithTransport() *HttpClientGo {
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
	}
	return client
}

func (client HttpClientGo) CheckDeletedDiag(d *schema.ResourceData, err error, response *http.Response, msg string) ([]byte, diag.Diagnostics) {
	if err != nil {
		return nil, diag.Errorf("%s: %s", msg, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, diag.Errorf("error convert data %s: %s", string(body), err)
	}

	if response.StatusCode == 200 {
		return body, nil
	}

	if strings.Contains(string(body), "does not exist") {
		resourceID := d.Id()
		d.SetId("")
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is goneand will be removed in Terraform state.", resourceID),
			},
		}
	}
	return nil, diag.Errorf("%s: %s", msg, err)
}
