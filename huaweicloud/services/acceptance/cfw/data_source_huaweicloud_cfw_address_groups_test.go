package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwAddressGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_address_groups.test"
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
				Config: testDataSourceDataSourceCfwAddressGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.object_id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.ref_count"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.address_type"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("key_word_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCfwAddressGroups_basic(name string) string {
	keyWord := "cfw_keyword"
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cfw_address_groups" "test" {
  depends_on  = [
    huaweicloud_cfw_address_group.test1,
    huaweicloud_cfw_address_group.test2,
    huaweicloud_cfw_address_group.test3,
  ]
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

data "huaweicloud_cfw_address_groups" "filter_by_name" {
  depends_on  = [
    huaweicloud_cfw_address_group.test1,
    huaweicloud_cfw_address_group.test2,
    huaweicloud_cfw_address_group.test3,
  ]
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = "%[2]s"
}

data "huaweicloud_cfw_address_groups" "filter_by_key_word" {
  depends_on  = [
    huaweicloud_cfw_address_group.test1,
    huaweicloud_cfw_address_group.test2,
    huaweicloud_cfw_address_group.test3,
  ]
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  key_word  = "%[3]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cfw_address_groups.filter_by_name.address_groups) > 0 && alltrue(
    [for v in data.huaweicloud_cfw_address_groups.filter_by_name.address_groups[*].name : v == "%[2]s"]
  )
}
	
output "key_word_filter_is_useful" {
  value = length(data.huaweicloud_cfw_address_groups.filter_by_key_word.address_groups) == 2
}
`, testAccDatasourceCreateAddressGroup(name, keyWord), name, keyWord)
}

func testAccDatasourceCreateAddressGroup(name, keyWord string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_address_group" "test1" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s"
  description = "%[3]s test"
}

resource "huaweicloud_cfw_address_group" "test2" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s_b"
  description = "%[3]s_d test"
}

resource "huaweicloud_cfw_address_group" "test3" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[2]s_c"
  description = "HTTP test"
}
`, testAccDatasourceFirewalls_basic(), name, keyWord)
}
