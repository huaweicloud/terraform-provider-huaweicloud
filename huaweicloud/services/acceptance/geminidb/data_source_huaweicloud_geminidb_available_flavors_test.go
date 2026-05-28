package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceAvailableFlavors_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_available_flavors.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailableFlavors_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "current_flavor.#"),
					resource.TestCheckResourceAttrSet(dataSource, "current_flavor.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "current_flavor.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "current_flavor.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "current_flavor.0.az_status.%"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testAccDataSourceAvailableFlavors_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceAvailableFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_geminidb_available_flavors" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testAccDataSourceAvailableFlavors_base(name))
}
