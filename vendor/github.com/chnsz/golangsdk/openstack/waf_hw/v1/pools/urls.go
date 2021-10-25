/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package pools

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("pool")
}

func resourceURL(c *golangsdk.ServiceClient, poolID string) string {
	return c.ServiceURL("pool", poolID)
}

func bindingURL(c *golangsdk.ServiceClient, poolID string) string {
	return c.ServiceURL("pool", poolID, "bindings")
}

func bindingResourceURL(c *golangsdk.ServiceClient, poolID, elbID string) string {
	return c.ServiceURL("pool", poolID, "bindings", elbID)
}
