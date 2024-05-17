package instances

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("instances", serverID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func modifyURL(c *golangsdk.ServiceClient, serverID, action string) string {
	return c.ServiceURL("instances", serverID, action)
}

func portModifiedURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "modify-port")
}

func secondsLevelMonitoringURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "monitoring-by-seconds/switch")
}

func backupPolicyURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "backups/policy")
}

func availabilityZoneURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "migrate")
}

func remarkURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "remark")
}

func replicaSetNameURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "replica-set/name")
}

func slowLogStatusURL(c *golangsdk.ServiceClient, instanceId string, status string) string {
	return c.ServiceURL("instances", instanceId, "slowlog-desensitization", status)
}

func restartURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "restart")
}

func maintenanceWindowURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "maintenance-window")
}

func balancerURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "balancer")
}

func balancerSwitchURL(c *golangsdk.ServiceClient, instanceId, action string) string {
	return c.ServiceURL("instances", instanceId, "balancer", action)
}

func balancerActiveWindowURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "balancer/active-window")
}

func clientNetworkRangesURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "client-network")
}
