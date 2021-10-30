package repositories

import (
	"github.com/chnsz/golangsdk"
)

const rootPath = "manage/namespaces"

func rootURL(client *golangsdk.ServiceClient, namespace string) string {
	return client.ServiceURL(rootPath, namespace, "repos")
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("manage", "repos")
}

func resourceURL(client *golangsdk.ServiceClient, namespace, repository string) string {
	return client.ServiceURL(rootPath, namespace, "repos", repository)
}
