package internal

import (
	"encoding/json"
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/response"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"time"
)

const (
	timeOut = 3
	host    = "169.254.169.254"
	errMsg  = "unable to get temporary credential"
	method  = "GET"
	path    = "/openstack/latest/securitykey"
)

var getTemporaryCredentialFromMetadataRequest = request.NewHttpRequestBuilder().
	WithEndpoint("http://" + host).
	WithMethod(method).
	WithPath(path).
	Build()

type GetTemporaryCredentialFromMetadataResponse struct {
	Credential *Credential `json:"credential,omitempty"`
}

type Credential struct {
	ExpiresAt string `json:"expires_at"`

	Access string `json:"access"`

	Secret string `json:"secret"`

	Securitytoken string `json:"securitytoken"`
}

func GetTemporaryCredential(client *impl.DefaultHttpClient) (*Credential, error) {

	type TempResp struct {
		value *response.DefaultHttpResponse
		err   error
	}

	respChan := make(chan TempResp, 1)

	go func() {
		defer close(respChan)

		resp, err := client.SyncInvokeHttp(getTemporaryCredentialFromMetadataRequest)
		respChan <- TempResp{
			value: resp,
			err:   err,
		}
	}()

	select {
	case tempResp := <-respChan:
		if tempResp.err != nil {
			return nil, tempResp.err
		}

		if tempResp.value.GetStatusCode() != 200 {
			return nil, sdkerr.NewServiceResponseError(tempResp.value.Response)
		}

		concreteResp := new(GetTemporaryCredentialFromMetadataResponse)
		err := json.Unmarshal([]byte(tempResp.value.GetBody()), concreteResp)

		if err != nil {
			return nil, err
		}

		if concreteResp.Credential == nil {
			return nil, errors.New(errMsg)
		}

		return concreteResp.Credential, nil
	case <-time.After(time.Second * timeOut):
		return nil, errors.New(errMsg)
	}

}
