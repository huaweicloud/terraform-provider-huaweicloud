/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package datamasking_rules

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, policy_id string) string {
	return c.ServiceURL("policy", policy_id, "privacy")
}

func resourceURL(c *golangsdk.ServiceClient, policy_id, id string) string {
	return c.ServiceURL("policy", policy_id, "privacy", id)
}
