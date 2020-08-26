package addons

import (
	"strings"

	"github.com/huaweicloud/golangsdk"
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

func CCEServiceURL(client *golangsdk.ServiceClient, cluster_id string, parts ...string) string {
	rbUrl := "https://" + cluster_id + "." + client.ResourceBaseURL()[8:]
	return rbUrl + strings.Join(parts, "/")
}
