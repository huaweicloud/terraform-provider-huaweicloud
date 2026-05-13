package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrimaryDimMonitoredObjects_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_primary_dim_monitored_objects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsPrimaryDimMonitoredObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "dim_name", "dcs_instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "router.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.dim_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.dim_route"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
				),
			},
		},
	})
}
func testAccDataSourceDcsPrimaryDimMonitoredObjects_basic() string {
	return `
data "huaweicloud_dcs_primary_dim_monitored_objects" "test" {
  dim_name = "dcs_instance_id"
}
`
}
