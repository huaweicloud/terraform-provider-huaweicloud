package block_devices

import "github.com/chnsz/golangsdk"

func getURL(c *golangsdk.ServiceClient, server_id string, volume_id string) string {
	return c.ServiceURL("cloudservers", server_id, "block_device", volume_id)
}

func rootURL(c *golangsdk.ServiceClient, serverId string) string {
	return c.ServiceURL("cloudservers", serverId, "block_device")
}

func resourceURL(c *golangsdk.ServiceClient, serverId, volumeId string) string {
	return c.ServiceURL("cloudservers", serverId, "block_device", volumeId)
}

func attachURL(c *golangsdk.ServiceClient, serverId string) string {
	return c.ServiceURL("cloudservers", serverId, "attachvolume")
}

func detachURL(c *golangsdk.ServiceClient, serverId, volumeId string) string {
	return c.ServiceURL("cloudservers", serverId, "detachvolume", volumeId)
}
