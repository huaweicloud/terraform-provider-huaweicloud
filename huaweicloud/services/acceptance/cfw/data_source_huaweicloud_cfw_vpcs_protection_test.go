package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwVpcsProtection_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_vpcs_protection.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwVpcsProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.total"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.total_assets"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.self_total"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.other_total"),
				),
			},
		},
	})
}

func testDataSourceCfwVpcsProtection_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}

locals {
  protect_objects = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects
  object_id       = try([for obj in local.protect_objects : obj.object_id if obj.type == 1][0], "")
}


data "huaweicloud_cfw_vpcs_protection" "test" {
  depends_on     = [data.huaweicloud_cfw_firewalls.test]
  object_id      = local.object_id
  fw_instance_id = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
