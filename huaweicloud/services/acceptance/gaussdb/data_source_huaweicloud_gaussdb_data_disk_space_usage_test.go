package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataDiskSpaceUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_data_disk_space_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataDiskSpaceUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_disk_capacity"),
					resource.TestCheckResourceAttrSet(dataSource, "data_disk_usage"),
					resource.TestCheckResourceAttrSet(dataSource, "space_usage_growth_per_day"),
					resource.TestCheckResourceAttrSet(dataSource, "estimated_remaining_days"),
					resource.TestCheckResourceAttrSet(dataSource, "cn_components.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dn_components.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dn_components.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dn_components.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dn_components.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "dn_components.0.node_name"),
				),
			},
		},
	})
}

func testDataSourceDataDiskSpaceUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_data_disk_space_usage" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
