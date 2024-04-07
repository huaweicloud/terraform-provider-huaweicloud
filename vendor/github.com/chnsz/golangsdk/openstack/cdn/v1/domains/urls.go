package domains

import (
	"fmt"
	"github.com/chnsz/golangsdk"
)

const (
	rootPath = "cdn/domains"
)

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(rootPath)
}

func updatePrivateBucketAccessURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "private-bucket-access")
}

func deleteURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId)
}

func getURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "detail")
}

func getDetailURL(sc *golangsdk.ServiceClient, domainName string) string {
	return sc.ServiceURL("cdn/configuration/domains", domainName)
}

func getTagsURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(fmt.Sprintf("cdn/configuration/tags?resource_id=%s", domainId))
}

func enableURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "enable")
}

func disableURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "disable")
}

func originURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "origin")
}
