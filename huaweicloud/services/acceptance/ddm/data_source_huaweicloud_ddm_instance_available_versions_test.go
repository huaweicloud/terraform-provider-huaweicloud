package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceDdmInstanceAvailableVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_instance_available_versions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdmInstanceAvailableVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "current_favored_version"),
					resource.TestCheckResourceAttrSet(dataSource, "current_version"),
					resource.TestCheckResourceAttrSet(dataSource, "latest_version"),
				),
			},
		},
	})
}

func testDataSourceDdmInstanceAvailableVersions_base(name string) string {
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

func testDataSourceDdmInstanceAvailableVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ddm_instance_available_versions" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
}
`, testDataSourceDdmInstanceAvailableVersions_base(name))
}
