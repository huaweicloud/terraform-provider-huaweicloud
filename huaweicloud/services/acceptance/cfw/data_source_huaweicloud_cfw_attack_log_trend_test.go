package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAttackLogTrend_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_attack_log_trend.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAttackLogTrend_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.deny_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.permit_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.time"),
				),
			},
		},
	})
}

func testDataSourceAttackLogTrend_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_log_trend" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  range          = "2"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
