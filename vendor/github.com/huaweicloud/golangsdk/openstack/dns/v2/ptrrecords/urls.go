package ptrrecords

import "github.com/huaweicloud/golangsdk"

func baseURL(c *golangsdk.ServiceClient, region string, floatingip_id string) string {
	return c.ServiceURL("reverse/floatingips", region+":"+floatingip_id)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("reverse/floatingips", id)
}
