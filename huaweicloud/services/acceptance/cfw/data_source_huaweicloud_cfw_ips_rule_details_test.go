package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwIpsRuleDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_ips_rule_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwIpsRuleDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ips_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ips_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceCfwIpsRuleDetails_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_ips_rule_details" "test" {
  fw_instance_id = "%[1]s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
