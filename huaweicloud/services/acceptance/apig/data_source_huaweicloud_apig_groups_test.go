package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGroupsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_groups.filter_by_name"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSource_filterByName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccGroupsDataSource_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccGroup_base(rName), rName)
}

func testAccGroupsDataSource_filterByName(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_groups" "filter_by_name" {
  // The behavior of parameter 'name' is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_apig_group.test,
  ]

  instance_id = huaweicloud_apig_instance.test.id
  name        = huaweicloud_apig_group.test.name
}

data "huaweicloud_apig_groups" "not_found" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_apig_group.test,
  ]  

  instance_id = huaweicloud_apig_instance.test.id
  name        = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_apig_groups.filter_by_name.groups[*].id : v == huaweicloud_apig_group.test.id]
}

output "is_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_apig_groups.not_found.groups) == 0
}
`, testAccGroupsDataSource_base(name))
}

func TestAccGroupsDataSource_filterById(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_groups.filter_by_id"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupsDataSource_filterById(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccGroupsDataSource_filterById(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_groups" "filter_by_id" {

  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id
}

data "huaweicloud_apig_groups" "not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_apig_group.test,
  ]  

  instance_id = huaweicloud_apig_instance.test.id
  group_id    = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_apig_groups.filter_by_id.groups[*].id : v == huaweicloud_apig_group.test.id]
}

output "is_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_apig_groups.not_found.groups) == 0
}
`, testAccGroupsDataSource_base(name), randUUID)
}
