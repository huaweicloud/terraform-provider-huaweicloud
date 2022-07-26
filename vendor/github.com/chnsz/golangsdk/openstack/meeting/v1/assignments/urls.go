package assignments

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/admin")
}

func resourceURL(client *golangsdk.ServiceClient, account string) string {
	return client.ServiceURL("usg/dcs/corp/admin", account)
}

func deleteURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/admin/delete")
}
