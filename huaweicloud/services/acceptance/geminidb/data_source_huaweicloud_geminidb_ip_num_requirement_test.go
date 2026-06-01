package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceIpNumRequirement_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_ip_num_requirement.test"
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
				Config: testAccDataSourceIpNumRequirement_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ip_address_count"),

					resource.TestCheckOutput("ip_address_number", "true"),
					resource.TestCheckOutput("instance_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceIpNumRequirement_base(name string) string {
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

func testAccDataSourceIpNumRequirement_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_geminidb_ip_num_requirement" "test" {
  node_num      = 2
  engine_name   = "redis"
  instance_mode = "Cluster"
}

data "huaweicloud_geminidb_ip_num_requirement" "instance_id_filter" {
  node_num    = 4
  instance_id = huaweicloud_geminidb_instance.test.id
}

output "ip_address_number" {
  value = data.huaweicloud_geminidb_ip_num_requirement.test.ip_address_count > 0
}

output "instance_id_filter_useful" {
  value = data.huaweicloud_geminidb_ip_num_requirement.instance_id_filter.ip_address_count > 0
}
`, testAccDataSourceIpNumRequirement_base(name))
}
