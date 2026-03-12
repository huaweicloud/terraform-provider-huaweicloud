package cse

import (
	"fmt"
	"io"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/auth"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
)

// getToken is a method to request the CSE API and get the authorization token.
// The format of is "Bearer {token}".
func getToken(c *golangsdk.ServiceClient, username, password string) (string, error) {
	tokenOpts := auth.CreateOpts{
		Name:     username,
		Password: password,
	}

	resp, err := auth.Create(c, tokenOpts)
	if err != nil {
		// When a microservice engine with an EIP is deleted, attempting to obtain a token via the engine's connection
		// address will result in an error (connection error). This needs to be handled specially to return a 404 error.
		if err == io.EOF || strings.Contains(err.Error(), "connection error") {
			return "", golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "POST",
					URL:       "/v4/token",
					RequestId: "NONE",
					Body:      []byte(`unable to connect the microservice engine`),
				},
			}
		}
		return "", fmt.Errorf("unable to create the authorization token: %v", err)
	}

	return fmt.Sprintf("Bearer %s", resp.Token), nil
}

// GetAuthorizationToken is a method for creating an authorization information for a CSE microservice to connect to the
// specified dedicated engine.
func GetAuthorizationToken(connAddr, username, password string) (string, error) {
	if username == "" {
		return "", nil
	}
	client := common.NewCustomClient(true, connAddr, "v4")
	return getToken(client, username, password)
}

func buildRequestMoreHeaders(enterpriseProjectId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
		"Accept":       "application/json",
	}

	if enterpriseProjectId != "" {
		moreHeaders["X-Enterprise-Project-ID"] = enterpriseProjectId
	}
	return moreHeaders
}
