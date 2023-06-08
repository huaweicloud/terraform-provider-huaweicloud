package cnad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceInstances_basic(t *testing.T) {
	rName := "data.huaweicloud_cnad_advanced_instances.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.region"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.instance_type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.protection_type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.ip_num"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.ip_num_now"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.protection_num"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.protection_num_now"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.created_at"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_type_filter_is_useful", "true"),
					resource.TestCheckOutput("region_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceInstances_basic() string {
	return `
data "huaweicloud_cnad_advanced_instances" "test" {
}

data "huaweicloud_cnad_advanced_instances" "instance_id_filter" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
}
output "instance_id_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_instances.instance_id_filter.instances[*].instance_id :
  v == data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id])
}

data "huaweicloud_cnad_advanced_instances" "instance_name_filter" {
  instance_name = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_name
}
output "instance_name_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_instances.instance_name_filter.instances[*].instance_name :
  v == data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_name])
}

data "huaweicloud_cnad_advanced_instances" "instance_type_filter" {
  instance_type = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_type
}
output "instance_type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_instances.instance_type_filter.instances[*].instance_type :
  v == data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_type])
}

data "huaweicloud_cnad_advanced_instances" "region_filter" {
  region = data.huaweicloud_cnad_advanced_instances.test.instances.0.region
}
output "region_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_cnad_advanced_instances.region_filter.instances[*].region :
  v == data.huaweicloud_cnad_advanced_instances.test.instances.0.region])
}
`
}
