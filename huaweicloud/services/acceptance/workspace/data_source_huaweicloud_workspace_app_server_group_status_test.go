package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerGroupStatus_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_app_server_group_status.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerGroupStatus_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "aps_status.%", regexp.MustCompile(`^[0-9]+$`)),
				),
			},
		},
	})
}

func testAccAppServerGroupStatus_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  app_type         = "COMMON_APP"
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"

  ip_virtual {
    enable = false
  }
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccAppServerGroupStatus_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_workspace_app_server_group_status" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
}
`, testAccAppServerGroupStatus_base(name))
}
