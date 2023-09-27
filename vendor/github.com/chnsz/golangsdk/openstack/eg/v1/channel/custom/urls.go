package custom

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("channels")
}

func resourceURL(c *golangsdk.ServiceClient, channelId string) string {
	return c.ServiceURL("channels", channelId)
}
