package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCpcsApps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpcs_apps.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCpcsApps_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.app_id"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.vpc_name"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.subnet_name"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.create_time"),

					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
					resource.TestCheckOutput("is_vpc_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCpcsApps_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cpcs_apps" "test" {
  depends_on = [huaweicloud_cpcs_app.test]
}

locals {
  app_name = data.huaweicloud_cpcs_apps.test.apps.0.app_name
  vpc_name = data.huaweicloud_cpcs_apps.test.apps.0.vpc_name
}

# Filter by app_name
data "huaweicloud_cpcs_apps" "app_name_filter" {
  app_name = local.app_name
}

locals {
  app_name_filter_result = [
    for v in data.huaweicloud_cpcs_apps.app_name_filter.apps[*].app_name : v == local.app_name
  ]
}

output "is_app_name_filter_useful" {
  value = length(local.app_name_filter_result) > 0 && alltrue(local.app_name_filter_result)
}

# Filter by vpc_name
data "huaweicloud_cpcs_apps" "vpc_name_filter" {
  vpc_name = local.vpc_name
}

locals {
  vpc_name_filter_result = [
    for v in data.huaweicloud_cpcs_apps.vpc_name_filter.apps[*].vpc_name : v == local.vpc_name
  ]
}

output "is_vpc_name_filter_useful" {
  value = length(local.vpc_name_filter_result) > 0 && alltrue(local.vpc_name_filter_result)
}
`, testCpcsApp_basic(name))
}
