package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwAttackLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_attack_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwAttackLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.packet"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.event_time"),
				),
			},
			{
				Config: testDataSourceCfwAttackLogs_app(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "records.0.app", "HTTP"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.packet"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.event_time"),
				),
			},
			{
				Config: testDataSourceCfwAttackLogs_level(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "records.0.level", "CRITICAL"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.packet"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.event_time"),
				),
			},
		},
	})
}

func testDataSourceCfwAttackLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_logs" "test" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_START_TIME, acceptance.HW_CFW_END_TIME)
}

func testDataSourceCfwAttackLogs_app() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_logs" "test" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  app            = "HTTP"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_START_TIME, acceptance.HW_CFW_END_TIME)
}

func testDataSourceCfwAttackLogs_level() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_logs" "test" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  level          = "CRITICAL"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_START_TIME, acceptance.HW_CFW_END_TIME)
}
