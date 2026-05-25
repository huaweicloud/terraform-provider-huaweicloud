package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccInstanceRestart_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testInstanceRestart_basic(name),
			},
		},
	})
}

func TestAccInstanceRestart_node_restart(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceRestart_node_restart(name),
			},
		},
	})
}

func testInstanceRestart_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_geminidb_flavors" "test" {
  engine_name  = "redis"
  mode         = "CloudNativeCluster"
  product_type = "Capacity"
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[1]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "CloudNativeCluster"
  product_type      = "Capacity"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_geminidb_flavors.test.flavors[1].spec_code
  }
}
`, common.TestBaseNetwork(name), name)
}

func testInstanceRestart_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_instance_restart" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testInstanceRestart_base(name))
}

func testAccInstanceRestart_node_restart(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_instance_restart" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  node_id     = huaweicloud_geminidb_instance.test.groups[0].nodes[0].id
}
`, testInstanceRestart_base(name))
}
