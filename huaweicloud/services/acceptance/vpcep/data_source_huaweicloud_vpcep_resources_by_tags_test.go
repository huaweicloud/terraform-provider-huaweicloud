package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcepResourcesByTags_basic(t *testing.T) {
	var (
		filterServicesByTagsAll  = "data.huaweicloud_vpcep_resources_by_tags.filter_services_by_tags"
		dcServicesTags           = acceptance.InitDataSourceCheck(filterServicesByTagsAll)
		filterEndpointsByTagsAll = "data.huaweicloud_vpcep_resources_by_tags.filter_endpoints_by_tags"
		dcEndpointsTags          = acceptance.InitDataSourceCheck(filterEndpointsByTagsAll)
		testName                 = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcepServicesByTags_basic(testName),
				Check: resource.ComposeTestCheckFunc(
					dcServicesTags.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(filterServicesByTagsAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterServicesByTagsAll, "resources.0.resource_id"),
					// 终端服务名称为region + name + 随机id拼接，暂不进行比较
					resource.TestCheckResourceAttrSet(filterServicesByTagsAll, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(filterServicesByTagsAll, "resources.0.tags.0.key", "foo0"),
					resource.TestCheckResourceAttr(filterServicesByTagsAll, "resources.0.tags.0.value", "bar0"),
					resource.TestCheckResourceAttr(filterServicesByTagsAll, "resources.0.tags.1.key", "foo1"),
					resource.TestCheckResourceAttr(filterServicesByTagsAll, "resources.0.tags.1.value", "bar1"),

					resource.TestCheckOutput("services_tags_filter_is_useful", "true"),
				),
			},
			{
				Config: testDataSourceVpcepEndpointsByTags_basic(testName),
				Check: resource.ComposeTestCheckFunc(
					dcEndpointsTags.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(filterEndpointsByTagsAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterEndpointsByTagsAll, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(filterEndpointsByTagsAll, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(filterEndpointsByTagsAll, "resources.0.tags.0.key", "foo4"),
					resource.TestCheckResourceAttr(filterEndpointsByTagsAll, "resources.0.tags.0.value", "bar4"),
					resource.TestCheckResourceAttr(filterEndpointsByTagsAll, "resources.0.tags.1.key", "foo5"),
					resource.TestCheckResourceAttr(filterEndpointsByTagsAll, "resources.0.tags.1.value", "bar5"),

					resource.TestCheckOutput("endpoints_tags_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testAccVPCResource_Precondition(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

resource "huaweicloud_compute_instance" "ecs" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpcep_service" "test" {
  name        = "%[2]s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    "foo0" = "bar0"
    "foo1" = "bar1"
  }
}
`, testAccCompute_data, rName)
}

func testDataSourceVpcepServicesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
data "huaweicloud_vpcep_resources_by_tags" "filter_services_by_tags" {
  depends_on = [huaweicloud_vpcep_service.test]
  action     = "filter"
  resource_type = "endpoint_service"

  tags {
    key    = "foo0"
    values = ["bar0"]
  }

  tags {
    key    = "foo1"
    values = ["bar1"]
  }
}

locals {
  tag_key   = "foo0"
  tag_value = "bar0"
}

output "services_tags_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_services_by_tags.resources) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_resources_by_tags.filter_services_by_tags.resources[*].tags : anytrue(
    [for vv in v[*].key : vv == local.tag_key]) && anytrue([for vv in v[*].value : vv == local.tag_value])]
  )
}

`, testAccVPCResource_Precondition(name))
}

func testDataSourceVpcepEndpointsByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_vpcep_endpoint" "test" {
  service_id       = huaweicloud_vpcep_service.test.id
  vpc_id           = data.huaweicloud_vpc.myvpc.id
  network_id       = data.huaweicloud_vpc_subnet.test.id
  enable_dns       = true
  description      = "test description"
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]

  tags = {
    "foo4" = "bar4"
    "foo5" = "bar5"
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_endpoints_by_tags" {
  depends_on = [huaweicloud_vpcep_endpoint.test]
  action     = "filter"
  resource_type = "endpoint"

  tags {
    key    = "foo4"
    values = ["bar4"]
  }

  tags {
    key    = "foo5"
    values = ["bar5"]
  }
}

locals {
  tag_key   = "foo4"
  tag_value = "bar4"
}

output "endpoints_tags_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_endpoints_by_tags.resources) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_resources_by_tags.filter_endpoints_by_tags.resources[*].tags : anytrue(
    [for vv in v[*].key : vv == local.tag_key]) && anytrue([for vv in v[*].value : vv == local.tag_value])]
  )
}

`, testAccVPCResource_Precondition(name))
}
