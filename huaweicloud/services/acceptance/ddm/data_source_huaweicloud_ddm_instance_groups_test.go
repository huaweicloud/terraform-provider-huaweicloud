package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceDdmInstanceGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	rName := "data.huaweicloud_ddm_instance_groups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmInstanceGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_list.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.id"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.name"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.role"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.endpoint"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.is_load_balance"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.is_default_group"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.cpu_num_per_node"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.mem_num_per_node"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.architecture"),
					resource.TestCheckResourceAttr(rName, "group_list.0.node_list.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.node_list.0.id"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.node_list.0.name"),
					resource.TestCheckResourceAttrSet(rName, "group_list.0.node_list.0.az"),
				),
			},
		},
	})
}

func testAccDatasourceDdmInstanceGroups_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceDdmInstanceGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ddm_instance_groups" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
}
`, testAccDatasourceDdmInstanceGroups_base(name))
}
