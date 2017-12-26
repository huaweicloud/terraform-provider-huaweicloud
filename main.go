package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/huawei-clouds/terraform-provider-huaweicloud/huaweicloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: huaweicloud.Provider})
}
