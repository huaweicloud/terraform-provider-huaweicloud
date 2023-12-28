package alarmreminding

import (
	"github.com/chnsz/golangsdk"
)

const resourcePathWarnAlert = "warnalert"
const resourcePathAlertConfig = "alertconfig"

func GetWarnAlertURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePathWarnAlert, resourcePathAlertConfig, "query")
}

func UpdateWarnAlertURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePathWarnAlert, resourcePathAlertConfig, "update")
}
