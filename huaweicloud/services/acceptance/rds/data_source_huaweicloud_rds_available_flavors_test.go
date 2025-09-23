package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsAvailableFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_available_flavors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsAvailableFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.is_ipv6_supported"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.type_code"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.az_status.%"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.max_connection"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.min_volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "optional_flavors.0.max_volume_size"),
					resource.TestCheckOutput("spec_code_like_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_category_type_filter_is_useful", "true"),
					resource.TestCheckOutput("is_rha_flavor_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsAvailableFlavors_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsAvailableFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_available_flavors" "test" {
  instance_id           = huaweicloud_rds_instance.test.id
  availability_zone_ids = "cn-north-4a"
  ha_mode               = "ha"
}

data "huaweicloud_rds_available_flavors" "spec_code_like_filter" {
  instance_id           = huaweicloud_rds_instance.test.id
  availability_zone_ids = "cn-north-4a"
  ha_mode               = "ha"
  spec_code_like        = data.huaweicloud_rds_available_flavors.test.optional_flavors[0].spec_code
}

output "spec_code_like_filter_is_useful" {
  value = length(data.huaweicloud_rds_available_flavors.spec_code_like_filter.optional_flavors) > 0
}

data "huaweicloud_rds_available_flavors" "flavor_category_type_filter" {
  instance_id           = huaweicloud_rds_instance.test.id
  availability_zone_ids = "cn-north-4a"
  ha_mode               = "ha"
  flavor_category_type  = "simple"
}

output "flavor_category_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_available_flavors.spec_code_like_filter.optional_flavors) > 0
}

data "huaweicloud_rds_available_flavors" "is_rha_flavor_filter" {
  instance_id           = huaweicloud_rds_instance.test.id
  availability_zone_ids = "cn-north-4a"
  ha_mode               = "ha"
  is_rha_flavor         = false
}

output "is_rha_flavor_filter_is_useful" {
  value = length(data.huaweicloud_rds_available_flavors.spec_code_like_filter.optional_flavors) > 0
}
`, testDataSourceRdsAvailableFlavors_base(name))
}
