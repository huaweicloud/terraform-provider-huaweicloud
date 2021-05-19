package credentials

import "github.com/huaweicloud/golangsdk"

const (
	rootPath        = "OS-CREDENTIAL"
	credentialsPath = "credentials"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, credentialsPath)
}

func resourceURL(client *golangsdk.ServiceClient, credID string) string {
	return client.ServiceURL(rootPath, credentialsPath, credID)
}
