package domains

import (
	"github.com/chnsz/golangsdk"
)

func rootURL(client *golangsdk.ServiceClient, namespace, repository string) string {
	return client.ServiceURL("manage/namespaces", namespace, "repositories", repository, "access-domains")
}

func resourceURL(client *golangsdk.ServiceClient, namespace, repository, domain string) string {
	return client.ServiceURL("manage/namespaces", namespace, "repositories", repository, "access-domains", domain)
}
