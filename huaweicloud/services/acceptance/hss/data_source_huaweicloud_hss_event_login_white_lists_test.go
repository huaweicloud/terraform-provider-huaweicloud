package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLoginWhiteLists_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_login_white_lists.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLoginWhiteLists_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.login_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.login_user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.remarks"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.enterprise_project_name"),

					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_login_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_login_user_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceLoginWhiteLists_base string = `
resource "huaweicloud_hss_event_login_white_list" "test" {
  private_ip            = "192.168.0.1"
  login_ip              = "192.168.0.2"
  login_user_name       = "user_test"
  remarks               = "remarks_test"
  handle_event          = true
  enterprise_project_id = "0"
}
`

func testAccDataSourceLoginWhiteLists_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_event_login_white_lists" "test" {
  depends_on = [huaweicloud_hss_event_login_white_list.test]
}

# Filter using private_ip.
locals {
  private_ip = data.huaweicloud_hss_event_login_white_lists.test.data_list[0].private_ip
}

data "huaweicloud_hss_event_login_white_lists" "private_ip_filter" {
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_login_white_lists.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_login_white_lists.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

# Filter using login_ip.
locals {
  login_ip = data.huaweicloud_hss_event_login_white_lists.test.data_list[0].login_ip
}

data "huaweicloud_hss_event_login_white_lists" "login_ip_filter" {
  login_ip = local.login_ip
}

output "is_login_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_login_white_lists.login_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_login_white_lists.login_ip_filter.data_list[*].login_ip : v == local.login_ip]
  )
}

# Filter using login_user_name.
locals {
  login_user_name = data.huaweicloud_hss_event_login_white_lists.test.data_list[0].login_user_name
}

data "huaweicloud_hss_event_login_white_lists" "login_user_name_filter" {
  login_user_name = local.login_user_name
}

output "is_login_user_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_login_white_lists.login_user_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_login_white_lists.login_user_name_filter.data_list[*].login_user_name : v == local.login_user_name]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_event_login_white_lists" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_hss_event_login_white_list.test]

  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_event_login_white_lists.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent login_user_name.
data "huaweicloud_hss_event_login_white_lists" "not_found" {
  login_user_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_event_login_white_lists.not_found.data_list) == 0
}
`, testDataSourceLoginWhiteLists_base)
}
