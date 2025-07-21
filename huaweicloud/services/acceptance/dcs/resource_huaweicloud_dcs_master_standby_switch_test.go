package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsMasterStandbySwitch_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsMasterStandbySwitch_basic(name),
			},
		},
	})
}

func testDcsMasterStandbySwitch_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
`, name)
}

func testDcsMasterStandbySwitch_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_master_standby_switch" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testDcsMasterStandbySwitch_base(name))
}
