package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopSysprep_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_sysprep.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDesktopSysprep_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "sysprep_info.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "sysprep_info.0.sysprep_version"),
					resource.TestCheckResourceAttrSet(dcName, "sysprep_info.0.support_create_image"),
				),
			},
		},
	})
}

func testAccDataDesktopSysprep_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_desktop_sysprep" "test" {
  desktop_id = huaweicloud_workspace_desktop.test.id
}
`, testAccDataDesktops_base(name))
}
