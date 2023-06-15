package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getProtectedInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return protectedinstances.Get(client, state.Primary.ID).Extract()
}

func TestAccProtectedInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_sdrs_protected_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectedInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectedInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "primary_ip_address", "192.168.0.15"),
					resource.TestCheckResourceAttr(rName, "delete_target_server", "true"),
					resource.TestCheckResourceAttr(rName, "delete_target_eip", "true"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "target_server"),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_sdrs_protection_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "server_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "primary_subnet_id", "huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testProtectedInstance_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cluster_id",
					"primary_subnet_id",
					"primary_ip_address",
					"delete_target_server",
					"delete_target_eip",
				},
			},
		},
	})
}

func testProtectedInstance_base(name string) string {
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
}
`, common.TestBaseNetwork(name), name)
}

func testProtectedInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

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
`, testProtectedInstance_base(name), name)
}

func testProtectedInstance_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_protected_instance" "test" {
  name                 = "%[2]s_update"
  group_id             = huaweicloud_sdrs_protection_group.test.id
  server_id            = huaweicloud_compute_instance.test.id
  primary_subnet_id    = huaweicloud_vpc_subnet.test.id
  primary_ip_address   = "192.168.0.15"
  delete_target_server = true
  delete_target_eip    = true
  description          = "test description"

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, testProtectedInstance_base(name), name)
}
