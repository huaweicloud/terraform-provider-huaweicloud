package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/checkpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getCheckpointResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	return checkpoints.Get(c, state.Primary.ID)
}

func TestAccCheckpoint_basic(t *testing.T) {
	var (
		policy       checkpoints.Checkpoint
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_cbr_checkpoint.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&policy,
		getCheckpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckpoint_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "backups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
				),
			},
		},
	})
}

func testAccCheckpoint_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = format("%[2]s-%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 10

  dynamic "resources" {
    for_each = huaweicloud_compute_instance.test[*].id

    content {
      server_id = resources.value
    }
  }
}

resource "huaweicloud_cbr_checkpoint" "test" {
  vault_id    = huaweicloud_cbr_vault.test.id
  name        = "%[2]s"
  description = "Created by terraform"

  dynamic "backups" {
    for_each = huaweicloud_compute_instance.test[*].id

    content {
      type        = "OS::Nova::Server"
      resource_id = backups.value
    }
  }
}
`, common.TestBaseComputeResources(name), name)
}
