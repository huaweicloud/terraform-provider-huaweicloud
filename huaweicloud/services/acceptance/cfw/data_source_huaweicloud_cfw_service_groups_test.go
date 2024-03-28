package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceServiceGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_service_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceServiceGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.ref_count"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "service_groups.0.protocols.#"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_not_found", "true"),
					resource.TestCheckOutput("key_word_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceServiceGroups_basic(rName string) string {
	keyWord := "cfw_sg_keyword"
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cfw_service_groups" "test" {
  depends_on  = [
    huaweicloud_cfw_service_group.test1,
    huaweicloud_cfw_service_group.test2,
    huaweicloud_cfw_service_group.test3,
  ]
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

data "huaweicloud_cfw_service_groups" "filter_by_name" {
  depends_on  = [
    huaweicloud_cfw_service_group.test1,
    huaweicloud_cfw_service_group.test2,
    huaweicloud_cfw_service_group.test3,
  ]
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = "%[2]s"
}

data "huaweicloud_cfw_service_groups" "filter_by_name_not_found" {
  depends_on  = [
    huaweicloud_cfw_service_group.test1,
    huaweicloud_cfw_service_group.test2,
    huaweicloud_cfw_service_group.test3,
  ]
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = "tf_test_not_found"
}

data "huaweicloud_cfw_service_groups" "filter_by_key_word" {
  depends_on  = [
    huaweicloud_cfw_service_group.test1,
    huaweicloud_cfw_service_group.test2,
    huaweicloud_cfw_service_group.test3,
  ]
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  key_word  = "%[3]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cfw_service_groups.filter_by_name.service_groups) > 0 && alltrue(
    [for v in data.huaweicloud_cfw_service_groups.filter_by_name.service_groups[*].name : v == "%[2]s"]
  )
}

output "name_filter_not_found" {
  value = length(data.huaweicloud_cfw_service_groups.filter_by_name_not_found.service_groups) == 0
}
	
output "key_word_filter_is_useful" {
  value = length(data.huaweicloud_cfw_service_groups.filter_by_key_word.service_groups) == 2
}
`, testAccDatasourceCreateServiceGroup(rName, keyWord), rName, keyWord)
}

func testAccDatasourceCreateServiceGroup(rName, keyWord string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_service_group" "test1" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s"
  description = "%[3]s"
}

resource "huaweicloud_cfw_service_group" "test2" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s_b"
  description = "%[3]s_a"
}

resource "huaweicloud_cfw_service_group" "test3" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s_c"
  description = "terraform test"
}
`, testAccDatasourceFirewalls_basic(), rName, keyWord)
}
