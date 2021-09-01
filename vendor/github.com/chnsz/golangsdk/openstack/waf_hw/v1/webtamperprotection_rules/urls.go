/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package webtamperprotection_rules

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL("policy", policyID, "antitamper")
}

func resourceURL(c *golangsdk.ServiceClient, policyID, id string) string {
	return c.ServiceURL("policy", policyID, "antitamper", id)
}
