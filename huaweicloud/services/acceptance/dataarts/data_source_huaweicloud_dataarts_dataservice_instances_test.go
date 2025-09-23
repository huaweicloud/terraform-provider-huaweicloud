package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Since the current service lacks resource to create an exclusive cluster, please ensure that at least one exclusive
// cluster has been created before testing.
func TestAccDataSourceDataServiceInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_dataservice_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_dataservice_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_dataarts_dataservice_instances.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byCreateUser   = "data.huaweicloud_dataarts_dataservice_instances.filter_by_create_user"
		dcByCreateUser = acceptance.InitDataSourceCheck(byCreateUser)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataServiceInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(all, "instances.0.id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "instances.0.create_user"),
					resource.TestCheckResourceAttrSet(all, "instances.0.status"),
					resource.TestMatchResourceAttr(all, "instances.0.flavor.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.flavor.0.id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.flavor.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.flavor.0.disk_size"),
					resource.TestCheckResourceAttrSet(all, "instances.0.flavor.0.vcpus"),
					resource.TestCheckResourceAttrSet(all, "instances.0.flavor.0.memory"),
					resource.TestCheckResourceAttrSet(all, "instances.0.gateway_version"),
					resource.TestCheckResourceAttrSet(all, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(all, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.security_group_id"),
					resource.TestMatchResourceAttr(all, "instances.0.node_num", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "instances.0.nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.private_ip"),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.status"),
					resource.TestMatchResourceAttr(all, "instances.0.nodes.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.create_user"),
					resource.TestCheckResourceAttrSet(all, "instances.0.nodes.0.gateway_version"),
					resource.TestMatchResourceAttr(all, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByCreateUser.CheckResourceExists(),
					resource.TestCheckOutput("is_create_user_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDataServiceInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_dataservice_instances" "test" {
  workspace_id = "%[1]s"
}

# Filter by name
locals {
  instance_name = data.huaweicloud_dataarts_dataservice_instances.test.instances[0].name
}

data "huaweicloud_dataarts_dataservice_instances" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.instance_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_instances.filter_by_name.instances[*].name : v == local.instance_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_dataarts_dataservice_instances" "filter_by_not_found_name" {
  workspace_id = "%[1]s"
  name         = local.not_found_name # This name is not exist 
}

locals {
  name_not_found_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_instances.filter_by_not_found_name.instances[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.name_not_found_filter_result) == 0
}

# Filter by create user
locals {
  create_user = data.huaweicloud_dataarts_dataservice_instances.test.instances[0].create_user
}

data "huaweicloud_dataarts_dataservice_instances" "filter_by_create_user" {
  workspace_id = "%[1]s"
  create_user  = local.create_user
}

locals {
  create_user_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_instances.filter_by_create_user.instances[*].create_user : v == local.create_user
  ]
}

output "is_create_user_filter_useful" {
  value = length(local.create_user_filter_result) > 0 && alltrue(local.create_user_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
