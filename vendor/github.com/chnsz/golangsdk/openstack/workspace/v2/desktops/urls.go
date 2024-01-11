package desktops

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("desktops")
}

func resourceURL(c *golangsdk.ServiceClient, desktopId string) string {
	return c.ServiceURL("desktops", desktopId)
}

func productURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("desktops/resize")
}

func volumeURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("volumes")
}

func volumeExpandURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("volumes/expand")
}

func rebuildURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("desktops/rebuild")
}

func eipsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("eips")
}

func bindEipURL(c *golangsdk.ServiceClient, bindType string) string {
	return c.ServiceURL("eips", bindType)
}

func networkURL(c *golangsdk.ServiceClient, desktopId string) string {
	return c.ServiceURL("desktops", desktopId, "networks")
}

func actionURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("desktops/action")
}
