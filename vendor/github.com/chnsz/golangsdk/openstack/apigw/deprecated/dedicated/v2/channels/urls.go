package channels

import (
	"github.com/chnsz/golangsdk"
)

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-channels")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-channels", chanId)
}

func membersURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-channels", chanId, "members")
}

func memberURL(c *golangsdk.ServiceClient, instanceId, chanId, memberId string) string {
	return c.ServiceURL("instances", instanceId, "vpc-channels", chanId, "members", memberId)
}
