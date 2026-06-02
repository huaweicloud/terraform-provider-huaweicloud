package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussDbDrConfigurationReset_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testGaussDbDrConfigurationReset_basic(rName),
			},
		},
	})
}

func testGaussDbDrConfigurationReset_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "test_1234"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testGaussDbDrConfigurationReset_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_dr_configuration_reset" "test" {
  instance_id        = huaweicloud_gaussdb_instance.test.id
  opposite_data_cidr = huaweicloud_vpc_subnet.test.cidr
}
`, testGaussDbDrConfigurationReset_base(rName))
}
