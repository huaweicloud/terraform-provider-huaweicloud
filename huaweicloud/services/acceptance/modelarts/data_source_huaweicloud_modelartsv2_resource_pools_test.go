package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please support at least one of resource pool.
func TestAccDataSourceV2ResourcePools_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelartsv2_resource_pools.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byWorkspaceId   = "data.huaweicloud_modelartsv2_resource_pools.by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		byStatus   = "data.huaweicloud_modelartsv2_resource_pools.by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2ResourcePools_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "resource_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestMatchResourceAttr(byStatus, "resource_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataSourceV2ResourcePools_basic string = `
data "huaweicloud_modelartsv2_resource_pools" "all" {}

data "huaweicloud_modelartsv2_resource_pools" "by_workspace_id" {
  workspace_id = "0"
}

output "is_workspace_id_filter_useful" {
  value = alltrue([for o in data.huaweicloud_modelartsv2_resource_pools.by_workspace_id.resource_pools: o.workspace_id == "0"])
}

data "huaweicloud_modelartsv2_resource_pools" "by_status" {
  status = "created"
}
`
