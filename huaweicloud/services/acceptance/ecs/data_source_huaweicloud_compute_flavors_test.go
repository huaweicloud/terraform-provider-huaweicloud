package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEcsFlavorsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_compute_flavors.this"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

const testAccEcsFlavorsDataSource_basic = `
data "huaweicloud_compute_flavors" "this" {
  performance_type = "normal"
  cpu_core_count   = 2
  memory_size      = 4
}
`
