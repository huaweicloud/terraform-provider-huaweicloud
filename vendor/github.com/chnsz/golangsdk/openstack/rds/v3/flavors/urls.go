package flavors

import "github.com/chnsz/golangsdk"

func listURL(sc *golangsdk.ServiceClient, databasename string) string {
	return sc.ServiceURL("flavors", databasename)
}
