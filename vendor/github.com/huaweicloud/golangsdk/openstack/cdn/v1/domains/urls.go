package domains

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "cdn/domains"
)

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(rootPath)
}

func deleteURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId)
}

func getURL(sc *golangsdk.ServiceClient, domainId string) string {
	return sc.ServiceURL(rootPath, domainId, "detail")
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
