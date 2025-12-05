package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopTags_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktop_tags.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataDesktopTags_basic_invalidDesktopId(),
				ExpectError: regexp.MustCompile(`The resource does not exist or is a resource of another project. DesktopId is not exist.`),
			},
			{
				Config: testAccDataDesktopTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestCheckResourceAttrSet(all, "tags.0.value"),
				),
			},
		},
	})
}

func testAccDataDesktopTags_basic_invalidDesktopId() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
# Filter by 'desktop_id' parameter and with invalid value.
data "huaweicloud_workspace_desktop_tags" "invalid_desktop_id" {
  desktop_id = "%[1]s"
}
`, randomId)
}

func testAccDataDesktopTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

locals {
  cpu_flavors = [for v in data.huaweicloud_workspace_flavors.test.flavors : v if v.is_gpu == false]
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id)
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group.0.id,
  ]

  nic {
    network_id = data.huaweicloud_workspace_service.test.network_ids[0]
  }

  name        = "%[2]s"
  user_name   = "user-%[2]s"
  user_email  = "terraform@example.com"
  user_group  = "administrators"
  delete_user = true

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
}
`, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID, name)
}

func testAccDataDesktopTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_desktop_tags" "all" {
  desktop_id = huaweicloud_workspace_desktop.test.id
}

locals {
  tag_result = try(data.huaweicloud_workspace_desktop_tags.all.tags[0], {})
}
`, testAccDataDesktopTags_base(name))
}
