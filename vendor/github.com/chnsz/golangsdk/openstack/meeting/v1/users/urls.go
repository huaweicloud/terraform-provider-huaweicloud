package users

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/member")
}

func resourceURL(client *golangsdk.ServiceClient, account string) string {
	return client.ServiceURL("usg/dcs/corp/member", account)
}

func deleteURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/member/delete")
}
