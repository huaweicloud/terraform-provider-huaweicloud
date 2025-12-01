package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopRemoteConsole_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktop_remote_console.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataDesktopRemoteConsole_basic_invalidDesktopId(),
				ExpectError: regexp.MustCompile(`The desktop does not exist.`),
			},
			{
				Config: testAccDataDesktopRemoteConsole_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "remote_console.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "remote_console.0.type"),
					resource.TestCheckResourceAttrSet(all, "remote_console.0.url"),
				),
			},
		},
	})
}

func testAccDataDesktopRemoteConsole_basic_invalidDesktopId() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
# Without any filter parameter but with invalid value.
data "huaweicloud_workspace_desktop_remote_console" "invalid_desktop_id" {
  desktop_id = "%[1]s"
}
`, randomId)
}

func testAccDataDesktopRemoteConsole_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_desktop_remote_console" "all" {
  desktop_id = huaweicloud_workspace_desktop.test.id

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}
`, testAccDataDesktops_base(name))
}
