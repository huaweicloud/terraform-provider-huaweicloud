package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccProtectedInstanceAddNIC_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProtectedInstanceAddNIC_basic(),
			},
		},
	})
}

func testProtectedInstanceAddNIC_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_sdrs_domain" "test" {}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "CentOS 7.6 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  target_availability_zone = data.huaweicloud_availability_zones.test.names[1]
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test.id
  description              = "test description"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      security_group_ids,
    ]
  }
}

resource "huaweicloud_sdrs_protected_instance" "test" {
  name                 = "%[2]s"
  group_id             = huaweicloud_sdrs_protection_group.test.id
  server_id            = huaweicloud_compute_instance.test.id
  primary_subnet_id    = huaweicloud_vpc_subnet.test.id
  primary_ip_address   = "192.168.0.15"
  delete_target_server = true
  delete_target_eip    = true
  description          = "test description"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_networking_secgroup" "test1" {
  name                 = "%[2]s_test1"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup" "test2" {
  name                 = "%[2]s_test2"
  delete_default_rules = true
}

resource "huaweicloud_sdrs_protected_instance_add_nic" "test" {
  protected_instance_id = huaweicloud_sdrs_protected_instance.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id

  security_groups {
    id = huaweicloud_networking_secgroup.test1.id
  }

  security_groups {
    id = huaweicloud_networking_secgroup.test2.id
  }
}
`, common.TestBaseNetwork(name), name)
}
