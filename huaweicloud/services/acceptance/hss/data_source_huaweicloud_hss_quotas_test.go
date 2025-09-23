package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used_status"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.shared_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.enterprise_project_name"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.tags.key", "value"),

					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_used_status_filter_useful", "true"),
					resource.TestCheckOutput("is_quota_id_filter_useful", "true"),
					resource.TestCheckOutput("is_charging_mode_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceQuotas_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_quotas" "test" {
  depends_on = [huaweicloud_hss_quota.test]
}

# Filter using category and category value is **container_resource**.
data "huaweicloud_hss_quotas" "category_filter" {
  category = "container_resource"
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.category_filter.quotas) == 0
}

# Filter using version.
locals {
  version = data.huaweicloud_hss_quotas.test.quotas[0].version
}

data "huaweicloud_hss_quotas" "version_filter" {
  version = local.version
}

output "is_version_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.version_filter.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_hss_quotas.version_filter.quotas[*].version : v == local.version]
  )
}

# Filter using status.
locals {
  status = data.huaweicloud_hss_quotas.test.quotas[0].status
}

data "huaweicloud_hss_quotas" "status_filter" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.status_filter.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_hss_quotas.status_filter.quotas[*].status : v == local.status]
  )
}

# Filter using used_status.
locals {
  used_status = data.huaweicloud_hss_quotas.test.quotas[0].used_status
}

data "huaweicloud_hss_quotas" "used_status_filter" {
  used_status = local.used_status
}

output "is_used_status_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.used_status_filter.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_hss_quotas.used_status_filter.quotas[*].used_status : v == local.used_status]
  )
}

# Filter using quota ID.
locals {
  quota_id = data.huaweicloud_hss_quotas.test.quotas[0].id
}

data "huaweicloud_hss_quotas" "quota_id_filter" {
  quota_id = local.quota_id
}

output "is_quota_id_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.quota_id_filter.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_hss_quotas.quota_id_filter.quotas[*].id : v == local.quota_id]
  )
}

# Filter using charging mode.
locals {
  charging_mode = data.huaweicloud_hss_quotas.test.quotas[0].charging_mode
}

data "huaweicloud_hss_quotas" "charging_mode_filter" {
  charging_mode = local.charging_mode
}

output "is_charging_mode_filter_useful" {
  value = length(data.huaweicloud_hss_quotas.charging_mode_filter.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_hss_quotas.charging_mode_filter.quotas[*].charging_mode : v == local.charging_mode]
  )
}

# Filter using non existent quota ID.
data "huaweicloud_hss_quotas" "not_found" {
  quota_id = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_quotas.not_found.quotas) == 0
}
`, testAccQuota_basic)
}
