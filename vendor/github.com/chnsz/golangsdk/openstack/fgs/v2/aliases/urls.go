package aliases

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, functionUrn string) string {
	return client.ServiceURL("fgs/functions", functionUrn, "aliases")
}

func resourceURL(client *golangsdk.ServiceClient, functionUrn, aliasName string) string {
	return client.ServiceURL("fgs/functions", functionUrn, "aliases", aliasName)
}
