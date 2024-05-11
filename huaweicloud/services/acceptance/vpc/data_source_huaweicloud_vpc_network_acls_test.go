package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkAclsDataSource_basic(t *testing.T) {
	var (
		name        = acceptance.RandomAccResourceNameWithDash()
		dataSource1 = "data.huaweicloud_vpc_network_acls.basic"
		dataSource2 = "data.huaweicloud_vpc_network_acls.filter_by_name"
		dataSource3 = "data.huaweicloud_vpc_network_acls.filter_by_id"
		dataSource4 = "data.huaweicloud_vpc_network_acls.filter_by_eps_id"
		dataSource5 = "data.huaweicloud_vpc_network_acls.filter_by_enabled"
		dataSource6 = "data.huaweicloud_vpc_network_acls.filter_by_status"
		dc1         = acceptance.InitDataSourceCheck(dataSource1)
		dc2         = acceptance.InitDataSourceCheck(dataSource2)
		dc3         = acceptance.InitDataSourceCheck(dataSource3)
		dc4         = acceptance.InitDataSourceCheck(dataSource4)
		dc5         = acceptance.InitDataSourceCheck(dataSource5)
		dc6         = acceptance.InitDataSourceCheck(dataSource6)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclsDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					dc5.CheckResourceExists(),
					dc6.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					resource.TestCheckOutput("is_enabled_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccNetworkAclsDataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_network_acls" "basic" {
  depends_on = [huaweicloud_vpc_network_acl.test]
}

data "huaweicloud_vpc_network_acls" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_vpc_network_acl.test]
}

data "huaweicloud_vpc_network_acls" "filter_by_id" {
  network_acl_id = huaweicloud_vpc_network_acl.test.id

  depends_on = [huaweicloud_vpc_network_acl.test]
}

data "huaweicloud_vpc_network_acls" "filter_by_eps_id" {
  enterprise_project_id = "0"

  depends_on = [huaweicloud_vpc_network_acl.test]
}

data "huaweicloud_vpc_network_acls" "filter_by_enabled" {
  enabled = "true"

  depends_on = [huaweicloud_vpc_network_acl.test]
}

data "huaweicloud_vpc_network_acls" "filter_by_status" {
  status = "INACTIVE"

  depends_on = [huaweicloud_vpc_network_acl.test]
}

locals {
  name_filter_result   = [for v in data.huaweicloud_vpc_network_acls.filter_by_name.network_acls[*].name : v == "%[2]s"]
  id_filter_result     = [
    for v in data.huaweicloud_vpc_network_acls.filter_by_name.network_acls[*].id : v == huaweicloud_vpc_network_acl.test.id
  ]
  eps_id_filter_result = [for v in data.huaweicloud_vpc_network_acls.filter_by_eps_id.network_acls[*].enterprise_project_id : v == "0"]
  enabled_filter_result = [for v in data.huaweicloud_vpc_network_acls.filter_by_enabled.network_acls[*].enabled : v == true]
  status_filter_result = [for v in data.huaweicloud_vpc_network_acls.filter_by_status.network_acls[*].status : v == "INACTIVE"]
  
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_network_acls.basic.network_acls) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_eps_id_filter_useful" {
  value = alltrue(local.eps_id_filter_result) && length(local.eps_id_filter_result) > 0
}

output "is_enabled_filter_useful" {
  value = alltrue(local.enabled_filter_result) && length(local.enabled_filter_result) > 0
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testAccNetworkAcl_basic(name), name)
}
