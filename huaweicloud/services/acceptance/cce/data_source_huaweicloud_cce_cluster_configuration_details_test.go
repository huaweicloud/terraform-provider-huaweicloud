package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceClusterConfigurationDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_cluster_configuration_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCceClusterConfigurationDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.%"),
					resource.TestCheckOutput("cluster_type_filter_is_useful", "true"),
					resource.TestCheckOutput("cluster_version_filter_is_useful", "true"),
					resource.TestCheckOutput("cluster_id_filter_is_useful", "true"),
					resource.TestCheckOutput("network_mode_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCceClusterConfigurationDetails_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_cluster_configuration_details" "test" {}

data "huaweicloud_cce_cluster_configuration_details" "cluster_type_filter" {
  cluster_type = "ARM64"
}
output "cluster_type_filter_is_useful" {
  value = length(data.huaweicloud_cce_cluster_configuration_details.cluster_type_filter.configurations) > 0
}

data "huaweicloud_cce_cluster_configuration_details" "cluster_version_filter" {
  cluster_version = "v1.33"
}
output "cluster_version_filter_is_useful" {
  value = length(data.huaweicloud_cce_cluster_configuration_details.cluster_version_filter.configurations) > 0
}

data "huaweicloud_cce_cluster_configuration_details" "cluster_id_filter" {
  cluster_id = "%s"
}
output "cluster_id_filter_is_useful" {
  value = length(data.huaweicloud_cce_cluster_configuration_details.cluster_id_filter.configurations) > 0
}

data "huaweicloud_cce_cluster_configuration_details" "network_mode_filter" {
  network_mode = "eni"
}
output "network_mode_filter_is_useful" {
  value = length(data.huaweicloud_cce_cluster_configuration_details.network_mode_filter.configurations) > 0
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
