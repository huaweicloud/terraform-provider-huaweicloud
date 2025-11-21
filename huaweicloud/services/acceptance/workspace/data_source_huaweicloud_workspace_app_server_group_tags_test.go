package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppServerGroupTags_basic(t *testing.T) {
	var (
		all               = "data.huaweicloud_workspace_app_server_group_tags.all"
		dcAll             = acceptance.InitDataSourceCheck(all)
		byServerGroupId   = "data.huaweicloud_workspace_app_server_group_tags.test"
		dcByServerGroupId = acceptance.InitDataSourceCheck(byServerGroupId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppServerGroupTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestMatchResourceAttr(all, "tags.0.values.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					dcByServerGroupId.CheckResourceExists(),
					resource.TestMatchResourceAttr(byServerGroupId, "tags.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(byServerGroupId, "tags.0.key"),
					resource.TestCheckResourceAttr(byServerGroupId, "tags.0.values.#", "1"),
					resource.TestCheckOutput("is_all_tags_include_test_tags", "true"),
				),
			},
		},
	})
}

func testAccDataAppServerGroupTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 90
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"

  tags = {
    owner = "terraform"
    foo   = "bar"
  }
}

data "huaweicloud_workspace_app_server_group_tags" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
}

data "huaweicloud_workspace_app_server_group_tags" "all" {
  depends_on = [huaweicloud_workspace_app_server_group.test]
}

output "is_all_tags_include_test_tags" {
  value = alltrue([for v in data.huaweicloud_workspace_app_server_group_tags.test.tags:
    length([for vv in data.huaweicloud_workspace_app_server_group_tags.all.tags:
      vv.key == v.key && contains(vv.values, v.values[0])
    ]) > 0
  ])
}
`, acceptance.RandomAccResourceName(),
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
	)
}
