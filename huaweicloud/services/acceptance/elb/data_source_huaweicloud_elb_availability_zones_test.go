package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAvailabilityZones_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.0.protocol.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.list.0.category"),
					resource.TestCheckOutput("public_border_group_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAvailabilityZones_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
 name           = "%[2]s"
 ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id

 availability_zone = [
   data.huaweicloud_availability_zones.test.names[0]
 ]
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceAvailabilityZones_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_availability_zones" "test" {}

locals {
  public_border_group = "center"
}

data "huaweicloud_elb_availability_zones" "public_border_group_filter" {
  public_border_group = "center"
}

output "public_border_group_filter_is_useful" {
  value = length(data.huaweicloud_elb_availability_zones.public_border_group_filter.availability_zones) > 0 && alltrue(
  [for v in data.huaweicloud_elb_availability_zones.public_border_group_filter.availability_zones[*].list :length(v) > 0 && alltrue(
  [for vv in v[*].public_border_group : vv == local.public_border_group]
  )]
  )
}

locals {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

data "huaweicloud_elb_availability_zones" "loadbalancer_id_filter" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_availability_zones.loadbalancer_id_filter.availability_zones) > 0 && alltrue(
  [for v in data.huaweicloud_elb_availability_zones.loadbalancer_id_filter.availability_zones[*].list : length(v) > 0]
  )
}
`, testAccDatasourceAvailabilityZones_base(name))
}
