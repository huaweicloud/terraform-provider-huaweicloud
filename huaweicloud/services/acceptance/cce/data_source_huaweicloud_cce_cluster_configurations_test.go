package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterConfigurationsDataSource_basic(t *testing.T) {
	datasourceName := "data.huaweicloud_cce_cluster_configurations.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testClusterConfigurations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testClusterConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_cluster_configurations" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_cluster_configurations.test.configurations) > 0
}
`, testAccCluster_basic(name))
}
