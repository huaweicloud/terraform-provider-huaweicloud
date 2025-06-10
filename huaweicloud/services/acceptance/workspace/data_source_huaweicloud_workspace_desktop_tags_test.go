package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDesktopTags_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_tags.filterByDesktopId"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDesktopTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_valid", "true"),
					resource.TestCheckOutput("is_key_valid", "true"),
					resource.TestCheckOutput("is_value_valid", "true"),
				),
			},
			{
				Config:      testAccDataSourceDesktopTags_ExpectError(),
				ExpectError: regexp.MustCompile(`The resource does not exist or is a resource of another project. DesktopId is not exist.`),
			},
		},
	})
}

func testAccDataSourceDesktopTags_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

# filter by desktop id
data "huaweicloud_workspace_desktop_tags" "filterByDesktopId" {
  desktop_id = huaweicloud_workspace_desktop.test.id
}

locals {
  tag_result = try(data.huaweicloud_workspace_desktop_tags.filterByDesktopId.tags[0], {})
}

output "is_tags_valid" {
  value = length(data.huaweicloud_workspace_desktop_tags.filterByDesktopId.tags) > 0
}

output "is_key_valid" {
  value = try(local.tag_result.key != null && local.tag_result.key != "", false)
}

output "is_value_valid" {
  value = try(local.tag_result.value != null, false)
}
`, testAccDataSourceDesktopTags_base(rName))
}

func testAccDataSourceDesktopTags_ExpectError() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
# expect error where id non exist
data "huaweicloud_workspace_desktop_tags" "expectErrorWhereIdNonExist" {
  desktop_id = "%[1]s"
}
`, randomId)
}

func testAccDataSourceDesktopTags_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}

# Ready for Desktop
locals {
  name              = "%[1]s"
  user_name         = "%[1]s-user"
  data_volume_sizes = [50, 70]
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = data.huaweicloud_workspace_flavors.test.flavors[0].id
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "")
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group[0].id,
    data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id,
  ]

  nic {
    network_id = data.huaweicloud_workspace_service.test.network_ids[0]
  }

  name       = local.name
  user_name  = local.user_name
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  email_notification = true
  delete_user        = true

  lifecycle {
    ignore_changes = [ name ]
  }
}
`, rName)
}
