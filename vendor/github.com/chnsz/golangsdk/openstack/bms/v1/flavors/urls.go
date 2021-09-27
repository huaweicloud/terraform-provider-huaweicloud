package flavors

import (
	"github.com/chnsz/golangsdk"
)

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("baremetalservers", "flavors")
}
