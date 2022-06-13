package auth

import "github.com/chnsz/golangsdk"

const rootPath = "usg/acs/auth/appauth"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
