package auth

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/acs/auth/account")
}

func validateURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/acs/token/validate")
}
