package cnad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceProtectedObjects_basic(t *testing.T) {
	rName := "data.huaweicloud_cnad_advanced_protected_objects.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceProtectedObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.id"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.type"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_name"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_version"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.region"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.status"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.block_threshold"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.clean_threshold"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("protected_object_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceProtectedObjects_basic() string {
	return `
data "huaweicloud_cnad_advanced_protected_objects" "test" {
}

data "huaweicloud_cnad_advanced_protected_objects" "instance_id_filter" {
  instance_id = data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.id
}
output "instance_id_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_protected_objects.instance_id_filter.protected_objects[*].
  instance_id : v == data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.id])
}

data "huaweicloud_cnad_advanced_protected_objects" "ip_address_filter" {
  ip_address = data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.ip_address
}
output "ip_address_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_protected_objects.ip_address_filter.protected_objects[*].
  ip_address : v == data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.ip_address])
}

data "huaweicloud_cnad_advanced_protected_objects" "protected_object_id_filter" {
  protected_object_id = data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.id
}
output "protected_object_id_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_protected_objects.protected_object_id_filter.
  protected_objects[*].id : v == data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.id])
}

data "huaweicloud_cnad_advanced_protected_objects" "type_filter" {
  type = data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.type
}
output "type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_protected_objects.type_filter.
  protected_objects[*].type : v == data.huaweicloud_cnad_advanced_protected_objects.test.protected_objects.0.type])
}
`
}
