package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCommonTaskStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_common_task_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCommonTaskStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "running_num"),
					resource.TestCheckResourceAttrSet(dataSource, "last_task_start_time"),
				),
			},
		},
	})
}

const testAccDataSourceCommonTaskStatistics_basic = `
data "huaweicloud_hss_common_task_statistics" "test" {
  task_type             = "cluster_scan"
  enterprise_project_id = "0"
}
`
