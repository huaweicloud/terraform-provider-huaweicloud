package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterVolumeUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_cluster_volume_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterVolumeUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.group"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.disk_capacity"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.disk_used"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_info_list.0.percentage"),
				),
			},
		},
	})
}

func testDataSourceClusterVolumeUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_cluster_volume_usage" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
