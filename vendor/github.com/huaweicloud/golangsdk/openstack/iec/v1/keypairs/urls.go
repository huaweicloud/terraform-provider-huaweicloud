package keypairs

import (
	"github.com/huaweicloud/golangsdk"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-keypairs")
}

func DeleteURL(c *golangsdk.ServiceClient, KeyPairName string) string {
	return c.ServiceURL("os-keypairs", KeyPairName)
}

func GetURL(c *golangsdk.ServiceClient, KeyPairName string) string {
	return c.ServiceURL("os-keypairs", KeyPairName)
}
