package namespaces

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
)

const rootPath = "namespaces"

func rootURL(client *golangsdk.ServiceClient, clusterID string) string {
	return addons.CCEServiceURL(client, clusterID, rootPath)
}

func resourceURL(client *golangsdk.ServiceClient, clusterID, name string) string {
	return addons.CCEServiceURL(client, clusterID, rootPath, name)
}
