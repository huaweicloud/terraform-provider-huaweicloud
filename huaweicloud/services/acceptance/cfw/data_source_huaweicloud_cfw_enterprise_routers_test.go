package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEnterpriseRouters_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_enterprise_routers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID.
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEnterpriseRouters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "er_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "er_list.0.er_id"),
					resource.TestCheckResourceAttrSet(dataSource, "er_list.0.name"),
				),
			},
		},
	})
}

func testDataSourceEnterpriseRouters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_enterprise_routers" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
