package cse

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const connectAddressRegex = `https://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`

var connectAddressRegexPattern = regexp.MustCompile(fmt.Sprintf(`^(%[1]s)?/?(%[1]s)/(.*)$`, connectAddressRegex))

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

// getToken is a method to request the CSE API and get the authorization token.
// The format of is "Bearer {token}".
func getToken(client *golangsdk.ServiceClient, username, password string) (string, error) {
	httpUrl := "v4/token"

	getPath := client.Endpoint + httpUrl

	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(client.ProjectID),
		JSONBody: map[string]interface{}{
			"name":     username,
			"password": password,
		},
	}

	requestResp, err := client.Request("POST", getPath, &requestOpts)
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
		return "", err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", utils.PathSearch("token", respBody, "").(string)), nil
}

// GetAuthorizationToken is a method for creating an authorization information when the microservice engine's RBAC
// authentication is enabled.
func GetAuthorizationToken(connAddr, username, password string) (string, error) {
	if username == "" {
		return "", nil
	}

	return getToken(common.NewCustomClient(true, connAddr), username, password)
}
