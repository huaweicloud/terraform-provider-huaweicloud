package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceClusterUpgradeInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_cluster_upgrade_info.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceClusterUpgradeInfo_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.release"),
					resource.TestCheckResourceAttrSet(dataSource, "spec.0.version_info.0.patch"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCceClusterUpgradeInfo_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_cluster_upgrade_info" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
}
`, testAccCceCluster_config(name))
}
