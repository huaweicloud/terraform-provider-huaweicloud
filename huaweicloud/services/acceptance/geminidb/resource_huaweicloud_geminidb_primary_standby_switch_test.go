package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccPrimaryStandbySwitch_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPrimaryStandbySwitch_basic(name),
			},
		},
	})
}

func testPrimaryStandbySwitch_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_geminidb_flavors" "test" {
  engine_name = "redis"
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Replication"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "2"
    size      = "4"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_geminidb_flavors.test.flavors[0].spec_code
  }
}
`, common.TestBaseNetwork(name), name)
}

func testPrimaryStandbySwitch_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_primary_standby_switch" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testPrimaryStandbySwitch_base(name))
}
