package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_tags.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByKey   = "data.huaweicloud_workspace_tags.filter_by_key"
		dcFilterByKey = acceptance.InitDataSourceCheck(filterByKey)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestMatchResourceAttr(all, "tags.0.values.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'key' parameter.
					dcFilterByKey.CheckResourceExists(),
					resource.TestCheckOutput("is_key_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = data.huaweicloud_workspace_flavors.test.flavors[0].id
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group[0].id,
    data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id,
  ]

  nic {
    network_id = data.huaweicloud_workspace_service.test.network_ids[0]
  }

  name       = "%[2]s"
  user_name  = "%[2]s-user"
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
    ignore_changes = [
      name
    ]
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID, name)
}

func testAccDataTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_tags" "all" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  key = "owner"
}

data "huaweicloud_workspace_tags" "filter_by_key" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]

  key = local.key
}

locals {
  tag_filter_result = [
    for v in data.huaweicloud_workspace_tags.filter_by_key.tags[*].key : v == local.key
  ]
}

output "is_key_filter_useful" {
  value = length(local.tag_filter_result) > 0 && alltrue(local.tag_filter_result)
}
`, testAccDataTags_base(name))
}
