package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTagsDataSource_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_tags.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcByKeyName = "data.huaweicloud_workspace_tags.filter_by_key"
		dcByKey     = acceptance.InitDataSourceCheck(dcByKeyName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTags(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.values.#"),
					dcByKey.CheckResourceExists(),
					resource.TestCheckOutput("is_key_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTags(rName string) string {
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
`, testAccDataSourceTags_base(rName))
}

func testAccDataSourceTags_base(rName string) string {
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

  name       = "%[1]s"
  user_name  = "%[1]s-user"
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
