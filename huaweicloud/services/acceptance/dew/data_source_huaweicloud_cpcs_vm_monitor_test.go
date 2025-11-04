package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVmMonitor_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_vm_monitor.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmMonitor_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.average"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.sum"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.variance"),
					resource.TestCheckResourceAttrSet(dataSource, "datapoints.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "metric_name_output"),
					resource.TestCheckResourceAttrSet(dataSource, "max"),
					resource.TestCheckResourceAttrSet(dataSource, "average"),
				),
			},
		},
	})
}

const testAccDataSourceVmMonitor_basic = `
data "huaweicloud_cpcs_vm_monitor" "test" {
  namespace   = "ECS"
  metric_name = "mem_util"
}`
