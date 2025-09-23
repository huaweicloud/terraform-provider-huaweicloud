package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_resource_groups.filter_by_EpsID"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "resource_groups.0.enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(dataSource, "resource_groups.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_groups.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_groups.0.created_at"),
					resource.TestCheckOutput("is_EpsID_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  EpsID = "%[2]s"
  name  = huaweicloud_ces_resource_group.test1.name
  id    = huaweicloud_ces_resource_group.test2.id
  type  = "TAG"
}

data "huaweicloud_ces_resource_groups" "filter_by_EpsID" {
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_ces_resource_group.test1,
    huaweicloud_ces_resource_group.test2
  ]
}

output "is_EpsID_filter_useful" {
  value = length(data.huaweicloud_ces_resource_groups.filter_by_EpsID) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_resource_groups.filter_by_EpsID.resource_groups[*] : item.enterprise_project_id == local.EpsID]
  )
}

data "huaweicloud_ces_resource_groups" "filter_by_name" {
  group_name = local.name
  
  depends_on = [
    huaweicloud_ces_resource_group.test1,
    huaweicloud_ces_resource_group.test2
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_ces_resource_groups.filter_by_name) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_resource_groups.filter_by_name.resource_groups[*] : item.group_name == local.name]
  )
}

data "huaweicloud_ces_resource_groups" "filter_by_id" {
  group_id = local.id

  depends_on = [
    huaweicloud_ces_resource_group.test1,
    huaweicloud_ces_resource_group.test2
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_ces_resource_groups.filter_by_id) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_resource_groups.filter_by_id.resource_groups[*] : item.group_id == local.id]
  )
}

data "huaweicloud_ces_resource_groups" "filter_by_type" {
  type = local.type

  depends_on = [
    huaweicloud_ces_resource_group.test1,
    huaweicloud_ces_resource_group.test2
  ]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_ces_resource_groups.filter_by_type) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_resource_groups.filter_by_type.resource_groups[*] : item.type == local.type]
  )
}
`, testDataSourceCesGroups_base(), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceCesGroups_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_ces_resource_group" "test1" {
  name = "%[1]s1"
  type = "TAG"
  tags = {
    key = "value"
  }
}

resource "huaweicloud_ces_resource_group" "test2" {
  name                  = "%[1]s2"
  type                  = "EPS"
  enterprise_project_id = "%[2]s"
  associated_eps_ids    = ["0"]
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
