package imagedata

import "github.com/huaweicloud/golangsdk"

const (
	rootPath   = "images"
	uploadPath = "file"
	stagePath  = "stage"
)

// `imageDataURL(c,i)` is the URL for the binary image data for the
// image identified by ID `i` in the service `c`.
func uploadURL(c *golangsdk.ServiceClient, imageID string) string {
	return c.ServiceURL(rootPath, imageID, uploadPath)
}

func stageURL(c *golangsdk.ServiceClient, imageID string) string {
	return c.ServiceURL(rootPath, imageID, stagePath)
}

func downloadURL(c *golangsdk.ServiceClient, imageID string) string {
	return uploadURL(c, imageID)
}
