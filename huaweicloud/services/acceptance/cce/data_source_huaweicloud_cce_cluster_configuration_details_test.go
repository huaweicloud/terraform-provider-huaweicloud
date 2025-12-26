package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHuaweiCloudCceClusterConfigurationDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_cluster_configuration_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHuaweiCloudCceClusterConfigurationDetails_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

const testDataSourceHuaweiCloudCceClusterConfigurationDetails_basic = `
data "huaweicloud_cce_cluster_configuration_details" "test" {}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_cluster_configuration_details.test.configurations) > 0
}
`
