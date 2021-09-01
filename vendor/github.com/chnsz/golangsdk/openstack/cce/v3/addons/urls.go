package addons

import (
	"net/url"
	"strings"

	"github.com/chnsz/golangsdk"
)

const (
	rootPath = "addons"
)

func rootURL(client *golangsdk.ServiceClient, cluster_id string) string {
	return CCEServiceURL(client, cluster_id, rootPath)
}

func resourceURL(client *golangsdk.ServiceClient, id, cluster_id string) string {
	return CCEServiceURL(client, cluster_id, rootPath, id+"?cluster_id="+cluster_id)
}

func resourceListURL(client *golangsdk.ServiceClient, cluster_id string) string {
	return CCEServiceURL(client, cluster_id, rootPath+"?cluster_id="+cluster_id)
}

func CCEServiceURL(client *golangsdk.ServiceClient, cluster_id string, parts ...string) string {
	u, _ := url.Parse(client.ResourceBaseURL())
	u.Host = cluster_id + "." + u.Host
	rbUrl := u.String()
	return rbUrl + strings.Join(parts, "/")
}
