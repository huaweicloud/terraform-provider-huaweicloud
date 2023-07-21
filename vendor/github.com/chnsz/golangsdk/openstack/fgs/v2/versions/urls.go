package versions

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, functionUrn string) string {
	return client.ServiceURL("fgs/functions", functionUrn, "versions")
}
