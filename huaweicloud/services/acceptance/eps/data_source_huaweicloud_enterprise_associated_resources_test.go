package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectAssociatedResources_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_associated_resources.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectAssociatedResources_basic_test(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "associated_resources.#"),
				),
			},
		},
	})
}

const testAccCompute_data = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_vpc" "test" {
  name = "test_eps_vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  depends_on = [
    huaweicloud_vpc.test
  ]

  name              = "test_eps_vpc_subnet"
  cidr              = "192.168.0.0/24"
  gateway_ip        = "192.168.0.1"
  vpc_id            = huaweicloud_vpc.test.id
  description       = "created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_identity_projects" "test" {
  name = "cn-north-4"
}
`

func testAccComputeInstance_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  depends_on = [
    huaweicloud_vpc_subnet.test
  ]

  name               = "test_eps"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data)
}

func testAccDataEnterpriseProjectAssociatedResources_basic_test() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_enterprise_project_associated_resources" "test" {
  depends_on = [
    huaweicloud_compute_instance.test
  ]

  resource_id   = huaweicloud_compute_instance.test.id
  project_id    = data.huaweicloud_identity_projects.test.projects[0].id
  region_id     = data.huaweicloud_identity_projects.test.name
  resource_type = "ecs"
}

output "huaweicloud_enterprise_project_associated_resources" {
  value = data.huaweicloud_enterprise_project_associated_resources.test
}
`, testAccComputeInstance_basic())
}
