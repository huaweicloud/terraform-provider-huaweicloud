package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceTmsResourceInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_resource_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsResourceInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "errors.#"),
					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTmsResourceInstances_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_identity_projects" "test" {
  name = "cn-north-4"
}

resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name           = "%[2]s_1"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_loadbalancer" "loadbalancer_2" {
  name           = "%[2]s_2"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
`, common.TestVpc(name), name)
}

func testDataSourceTmsResourceInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_tms_resource_instances" "test" {
  depends_on = [
    huaweicloud_elb_loadbalancer.loadbalancer_1,
    huaweicloud_elb_loadbalancer.loadbalancer_2
  ]

  project_id     = data.huaweicloud_identity_projects.test.projects[0].id
  resource_types = ["loadbalancers"]

  tags {
    key    = "owner"
    values = ["terraform"]
  }
}

locals {
  without_any_tag = "true"
}
data "huaweicloud_tms_resource_instances" "without_any_tag_filter" {
  depends_on = [
    huaweicloud_elb_loadbalancer.loadbalancer_1,
    huaweicloud_elb_loadbalancer.loadbalancer_2
  ]

  project_id      = data.huaweicloud_identity_projects.test.projects[0].id
  resource_types  = ["loadbalancers"]
  without_any_tag = "true"

  tags {
    key    = "key"
    values = ["value"]
  }
}
output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_tms_resource_instances.without_any_tag_filter.resources) > 0 && alltrue(
  [for v in data.huaweicloud_tms_resource_instances.without_any_tag_filter.resources[*].tags : length(v) == 0]
  )  
}
`, testDataSourceTmsResourceInstances_base(name))
}
