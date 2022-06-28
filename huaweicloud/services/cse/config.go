package cse

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/auth"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
)

// getAuthorizationToken is a method to request the CSE API and get the authorization token.
// The format of is "Bearer {token}".
func getAuthorizationToken(c *golangsdk.ServiceClient, username, password string) (string, error) {
	tokenOpts := auth.CreateOpts{
		Name:     username,
		Password: password,
	}
	resp, err := auth.Create(c, tokenOpts)
	if err != nil {
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
	return getAuthorizationToken(client, username, password)
}
