package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostGroups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_host_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.host_ids.#"),

					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					resource.TestCheckOutput("is_host_num_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceHostGroups_basic() string {
	name := acceptance.RandomAccResourceName()
	hostGroupBasic := testAccHostGroup_basic(name)

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_hss_host_groups" "test" {
  depends_on = [huaweicloud_hss_host_group.test]

  enterprise_project_id = "%[2]s"
}

# Filter using group ID.
locals {
  group_id = data.huaweicloud_hss_host_groups.test.groups[0].id
}

data "huaweicloud_hss_host_groups" "group_id_filter" {
  group_id              = local.group_id
  enterprise_project_id = "%[2]s"
}

output "is_group_id_filter_useful" {
  value = length(data.huaweicloud_hss_host_groups.group_id_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_groups.group_id_filter.groups[*].id : v == local.group_id]
  )
}

# Filter using host_num.
locals {
  host_num = data.huaweicloud_hss_host_groups.test.groups[0].host_num
}

data "huaweicloud_hss_host_groups" "host_num_filter" {
  host_num              = local.host_num
  enterprise_project_id = "%[2]s"
}

output "is_host_num_filter_useful" {
  value = length(data.huaweicloud_hss_host_groups.host_num_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_groups.host_num_filter.groups[*].host_num : v == local.host_num]
  )
}

# Filter using non existent name.
data "huaweicloud_hss_host_groups" "not_found" {
  name                  = "resource_not_found"
  enterprise_project_id = "%[2]s"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_host_groups.not_found.groups) == 0
}
`, hostGroupBasic, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
