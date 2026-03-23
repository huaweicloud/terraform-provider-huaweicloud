package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAttackLogStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_attack_log_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAttackLogStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.app_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.attack_rule_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.attack_type_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dst_ip_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dst_port_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.src_ip_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
				),
			},
		},
	})
}

func testDataSourceAttackLogStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_log_statistics" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  action         = 0
  item           = "src_region_id"
  value          = "target-object"
  range          = "2"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
