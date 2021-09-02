/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package premium_instances

import "github.com/chnsz/golangsdk"

const (
	resourcePath = "instance"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
