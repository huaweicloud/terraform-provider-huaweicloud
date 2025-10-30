package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCpcsClusterAccessKeys_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpcs_cluster_access_keys.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCpcsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCpcsClusterAccessKeys_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.access_key_id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.create_time"),

					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCpcsClusterAccessKeys_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cpcs_cluster_access_keys" "test" {
  cluster_id = "%[1]s"
}

locals {
  app_name = data.huaweicloud_cpcs_cluster_access_keys.test.access_keys.0.app_name
}

# Filter by app name
data "huaweicloud_cpcs_cluster_access_keys" "app_name_filter" {
  cluster_id = "%[1]s"
  app_name   = local.app_name
}

locals {
  app_name_filter_result = [
    for v in data.huaweicloud_cpcs_cluster_access_keys.app_name_filter.access_keys[*].app_name : v == local.app_name
  ]
}

output "is_app_name_filter_useful" {
  value = length(local.app_name_filter_result) > 0 && alltrue(local.app_name_filter_result)
}
`, acceptance.HW_CPCS_CLUSTER_ID)
}
