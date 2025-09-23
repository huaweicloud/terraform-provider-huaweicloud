package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDdmInstanceRollback_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ddm_instance_rollback.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdmInstanceRollback_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDdmInstanceRollback_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.9"
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

data "huaweicloud_ddm_instance_available_versions" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
}

resource "huaweicloud_ddm_instance_upgrade" "test" {
  depends_on = [data.huaweicloud_ddm_instance_available_versions.test]

  instance_id    = huaweicloud_ddm_instance.test.id
  target_version = data.huaweicloud_ddm_instance_available_versions.test.latest_version
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDdmInstanceRollback_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_instance_rollback" "test" {
  depends_on = [huaweicloud_ddm_instance_upgrade.test]

  instance_id = huaweicloud_ddm_instance.test.id
}`, testAccDdmInstanceRollback_base(rName))
}
