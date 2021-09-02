/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package valuelists

import "github.com/chnsz/golangsdk"

const rootPath = "valuelist"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}
