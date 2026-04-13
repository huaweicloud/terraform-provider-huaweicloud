package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwFirewallRouter_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_firewall_router.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwFirewallRouter_basic(),
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

func testDataSourceCfwFirewallRouter_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewall_router" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
