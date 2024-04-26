package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcPermissions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_permissions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCPermission(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcPermissions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.cloud_connection_id"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("conn_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcPermissions_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cc_permissions" "test" {
  permission_id = "%[1]s"
}

locals {
  permissions         = data.huaweicloud_cc_permissions.test.permissions
  name                = local.permissions[0].name
  description         = local.permissions[0].description
  cloud_connection_id = local.permissions[0].cloud_connection_id
  instance_id         = local.permissions[0].instance_id
}

data "huaweicloud_cc_permissions" "filter_by_name" {
  name = local.name
}

data "huaweicloud_cc_permissions" "filter_by_name_not_found" {
  name = "%[2]s"
}

data "huaweicloud_cc_permissions" "filter_by_description" {
  description = local.description
}

data "huaweicloud_cc_permissions" "filter_by_desc_not_found" {
  description = "%[2]s"
}

data "huaweicloud_cc_permissions" "filter_by_conn_id" {
  cloud_connection_id = local.cloud_connection_id
}

data "huaweicloud_cc_permissions" "filter_by_instance_id" {
  instance_id = local.instance_id
}

locals {
  list_by_name           = data.huaweicloud_cc_permissions.filter_by_name.permissions
  list_by_name_not_found = data.huaweicloud_cc_permissions.filter_by_name_not_found.permissions
  list_by_description    = data.huaweicloud_cc_permissions.filter_by_description.permissions
  list_by_desc_not_found = data.huaweicloud_cc_permissions.filter_by_desc_not_found.permissions
  list_by_conn_id        = data.huaweicloud_cc_permissions.filter_by_conn_id.permissions
  list_by_instence_id    = data.huaweicloud_cc_permissions.filter_by_instance_id.permissions
}

output "id_filter_is_useful" {
  value = length(local.permissions) > 0 && alltrue(
    [for v in local.permissions[*].id : v == "%[1]s"]
  )
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && length(local.list_by_name_not_found) == 0
}

output "description_filter_is_useful" {
  value = length(local.list_by_description) > 0 && length(local.list_by_desc_not_found) == 0
}

output "conn_id_filter_is_useful" {
  value = length(local.list_by_conn_id) > 0 && alltrue(
    [for v in local.list_by_conn_id[*].cloud_connection_id : v == local.cloud_connection_id]
  )
}

output "instance_id_filter_is_useful" {
  value = length(local.list_by_instence_id) > 0 && alltrue(
    [for v in local.list_by_instence_id[*].instance_id : v == local.instance_id]
  )
}
`, acceptance.HW_CC_PERMISSION_ID, name)
}
