package antiddos

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId)
}

func DailyReportURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId, "daily")
}

func DeleteURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId)
}

func GetURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId)
}

func GetStatusURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId, "status")
}

func GetTaskURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("query_task_status")
}

func ListConfigsURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("antiddos", "query_config_list")
}

func ListLogsURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId, "logs")
}

func ListStatusURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("antiddos")
}

func UpdateURL(c *golangsdk.ServiceClient, floatingIpId string) string {
	return c.ServiceURL("antiddos", floatingIpId)
}

func WeeklyReportURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("antiddos", "weekly")
}
