package cloudimages

import (
	"net/url"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/utils"
)

// query images using search criteria and to display the images in a list
func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudimages")
}

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudimages/action")
}

func createDataImageURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudimages/dataimages/action")
}

// builds next page full url based on current url
func nextPageURL(serviceURL, requestedNext string) (string, error) {
	base, err := utils.BaseEndpoint(serviceURL)
	if err != nil {
		return "", err
	}

	requestedNextURL, err := url.Parse(requestedNext)
	if err != nil {
		return "", err
	}

	base = golangsdk.NormalizeURL(base)
	nextPath := base + strings.TrimPrefix(requestedNextURL.Path, "/")

	nextURL, err := url.Parse(nextPath)
	if err != nil {
		return "", err
	}

	nextURL.RawQuery = requestedNextURL.RawQuery
	return nextURL.String(), nil
}
