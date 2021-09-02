package templates

import "github.com/chnsz/golangsdk"

func List(client *golangsdk.ServiceClient, cluster_id string) (r ListResutlt) {
	_, r.Err = client.Get(templateURL(client, cluster_id), &r.Body, nil)
	return
}
