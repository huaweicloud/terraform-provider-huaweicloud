package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopVolumeBatchDelete_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		desktop      interface{}
		resourceName = "huaweicloud_workspace_desktop.test"
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&desktop,
			getDesktopFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopVolumeBatchDelete_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckOutput("is_desktop_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDesktopVolumeBatchDelete_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  image_id = "%[2]s"
  visibility = "market"
}

variable "data_volumes" {
  type = list(object({
    type       = string
    size       = number
  }))

  default = [
    { type = "SSD", size = 80 },
    { type = "SSD", size = 100 },
  ]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"), 
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[1]s"
  user_name  = "user-%[1]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  dynamic "data_volume" {
    for_each = var.data_volumes

    content {
      type = data_volume.value.type
      size = data_volume.value.size
    }
  }

  lifecycle {
    ignore_changes = ["data_volume"]
  }
}
`, name, acceptance.HW_IMAGE_ID)
}

func testAccDesktopVolumeBatchDelete_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop_volume_batch_delete" "test" {
  desktop_id = huaweicloud_workspace_desktop.test.id
  volume_ids = [
    huaweicloud_workspace_desktop.test.data_volume.0.id
  ]

  lifecycle {
    ignore_changes = ["volume_ids"]
  }

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

# By desktop ID filter
data "huaweicloud_workspace_desktops" "filter_by_desktop_id" {
  desktop_id = huaweicloud_workspace_desktop.test.id

  depends_on = [
    huaweicloud_workspace_desktop_volume_batch_delete.test
  ]
}

output "is_desktop_id_filter_useful" {
  value = (
    length(data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops) == 1 &&
    length(data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops[0].data_volume) == 1
  )
}
`, testAccDesktopVolumeBatchDelete_base(name))
}
