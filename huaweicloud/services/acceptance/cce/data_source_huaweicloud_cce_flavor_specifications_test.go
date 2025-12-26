package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHuaweiCloudCceFlavorSpecifications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_flavor_specifications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHuaweiCloudCceFlavorSpecifications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

const testDataSourceHuaweiCloudCceFlavorSpecifications_basic = `
data "huaweicloud_cce_flavor_specifications" "test" {
  cluster_type = "VirtualMachine"
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_flavor_specifications.test.cluster_flavor_specs) > 0
}
`
