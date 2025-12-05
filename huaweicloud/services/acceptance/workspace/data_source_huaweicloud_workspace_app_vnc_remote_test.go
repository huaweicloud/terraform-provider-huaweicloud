package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppVncRemote_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_app_vnc_remote.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppVncRemote_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "url", regexp.MustCompile(`^https?://.*$`)),
					resource.TestCheckOutput("is_type_valid", "true"),
				),
			},
		},
	})
}

func testAccDataAppVncRemote_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_vnc_remote" "test" {
  server_id = "%[1]s"
}

output "is_type_valid" {
  value = strcontains(data.huaweicloud_workspace_app_vnc_remote.test.type, "vnc")
}
`, acceptance.HW_WORKSPACE_APP_SERVER_ID)
}
