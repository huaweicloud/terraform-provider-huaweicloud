package templates

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
)

const templatePath = "addontemplates"

func templateURL(client *golangsdk.ServiceClient, cluster_id string) string {
	return addons.CCEServiceURL(client, cluster_id, templatePath)
}
