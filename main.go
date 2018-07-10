package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-huaweicloud/huaweicloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: huaweicloud.Provider})
}
