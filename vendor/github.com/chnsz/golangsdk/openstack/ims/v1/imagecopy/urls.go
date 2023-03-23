package imagecopy

import "github.com/chnsz/golangsdk"

func withinRegionCopyURL(c *golangsdk.ServiceClient, imageId string) string {
	return c.ServiceURL("cloudimages", imageId, "copy")
}

func crossRegionCopyURL(c *golangsdk.ServiceClient, imageId string) string {
	return c.ServiceURL("cloudimages", imageId, "cross_region_copy")
}
