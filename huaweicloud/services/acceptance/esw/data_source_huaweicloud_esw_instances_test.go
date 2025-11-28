package esw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_esw_instances.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instances.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.project_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.region"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.flavor_ref"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.ha_mode"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.description"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.availability_zones.0.primary"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.availability_zones.0.standby"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.0.virsubnet_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.0.tunnel_ip"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.0.tunnel_port"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tunnel_info.0.tunnel_type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.charge_infos.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.charge_infos.0.charge_mode"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceEswInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_esw_instances" "test" {
  depends_on = [huaweicloud_esw_instance.test]
}

locals{
  instance_id = huaweicloud_esw_instance.test.id
}
data "huaweicloud_esw_instances" "instance_id_filter" {
  depends_on = [huaweicloud_esw_instance.test]

  instance_id = huaweicloud_esw_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_instances.instance_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_esw_instances.instance_id_filter.instances[*].id : v == local.instance_id]
  )
}

locals{
  name = huaweicloud_esw_instance.test.name
}
data "huaweicloud_esw_instances" "name_filter" {
  depends_on = [huaweicloud_esw_instance.test]

  name = huaweicloud_esw_instance.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_esw_instances.name_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_esw_instances.name_filter.instances[*].name : v == local.name]
  )
}

locals{
  description = huaweicloud_esw_instance.test.description
}
data "huaweicloud_esw_instances" "description_filter" {
  depends_on = [huaweicloud_esw_instance.test]

  description = huaweicloud_esw_instance.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_esw_instances.description_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_esw_instances.description_filter.instances[*].description : v == local.description]
  )
}
`, testAccEswInstance_basic(name))
}
