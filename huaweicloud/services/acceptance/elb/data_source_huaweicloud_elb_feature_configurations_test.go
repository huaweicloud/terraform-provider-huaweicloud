package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbFeatureConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_feature_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbFeatureConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.feature"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.switch"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.caller"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.updated_at"),
					resource.TestCheckOutput("feature_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbFeatureConfigurations_basic() string {
	return `
data "huaweicloud_elb_feature_configurations" "test" {}

locals {
  feature = data.huaweicloud_elb_feature_configurations.test.configs[0].feature
}

data "huaweicloud_elb_feature_configurations" "feature_filter" {
  feature = data.huaweicloud_elb_feature_configurations.test.configs[0].feature
}

output "feature_filter_is_useful" {
  value = length(data.huaweicloud_elb_feature_configurations.feature_filter.configs) > 0 && alltrue(
  [for v in data.huaweicloud_elb_feature_configurations.feature_filter.configs[*].feature : v == local.feature]
  )
}
`
}
