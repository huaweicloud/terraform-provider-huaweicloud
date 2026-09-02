package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafPools_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_waf_pools.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWafPools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.create_time"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceWafPools_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_waf_pool" "test" {
  name                  = "%[1]s"
  type                  = "detector-cloud"
  vpc_id                = huaweicloud_vpc.test.id
  description           = "created by terraform"
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceWafPools_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_pools" "test" {
  depends_on = [huaweicloud_waf_pool.test]

  enterprise_project_id = "%[2]s"
  detail                = true
}

# Filter by name.
locals {
  name = data.huaweicloud_waf_pools.test.items[0].name
}

data "huaweicloud_waf_pools" "filter_by_name" {
  enterprise_project_id = "%[2]s"
  name                  = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_pools.filter_by_name.items[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by type.
locals {
  type = data.huaweicloud_waf_pools.test.items[0].type
}

data "huaweicloud_waf_pools" "filter_by_type" {
  enterprise_project_id = "%[2]s"
  type                  = [local.type]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_waf_pools.filter_by_type.items[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

# Filter by vpc_id.
locals {
  vpc_id = data.huaweicloud_waf_pools.test.items[0].vpc_id
}

data "huaweicloud_waf_pools" "filter_by_vpc_id" {
  enterprise_project_id = "%[2]s"
  vpc_id                = local.vpc_id
}

locals {
  vpc_id_filter_result = [
    for v in data.huaweicloud_waf_pools.filter_by_vpc_id.items[*].vpc_id : v == local.vpc_id
  ]
}

output "vpc_id_filter_is_useful" {
  value = alltrue(local.vpc_id_filter_result) && length(local.vpc_id_filter_result) > 0
}
`, testDataSourceWafPools_base(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
