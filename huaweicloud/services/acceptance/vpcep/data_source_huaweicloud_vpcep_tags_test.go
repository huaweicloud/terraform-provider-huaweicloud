package vpcep

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		all   = "data.huaweicloud_vpcep_tags.test"
		dc    = acceptance.InitDataSourceCheck(all)
		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestMatchResourceAttr(all, "tags.0.values.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("tags_validation", "true"),
				),
			},
		},
	})
}

func testAccDataTags_base(rName string) string {
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
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }

  tags = {
    owner = "tf-acc"
  }
}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id = huaweicloud_vpcep_service.test.id
  vpc_id     = data.huaweicloud_vpc.myvpc.id
  network_id = data.huaweicloud_vpc_subnet.test.id
  enable_dns = true

  tags = {
    owner = "tf-acc"
  }

  lifecycle {
    ignore_changes = [enable_dns]
  }
}

`, testAccCompute_data, rName, rName)
}

func testAccDataTags_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_tags" "test" {
  depends_on    = [huaweicloud_vpcep_endpoint.test]
  resource_type = "endpoint"
}

output "tags_validation" {
  value = length([for t in data.huaweicloud_vpcep_tags.test.tags: t.key == "owner" &&
    alltrue([for k, v in huaweicloud_vpcep_endpoint.test[*].tags: contains(t.values, v) if k == "owner"])]) > 0
}

`, testAccDataTags_base(rName))
}
