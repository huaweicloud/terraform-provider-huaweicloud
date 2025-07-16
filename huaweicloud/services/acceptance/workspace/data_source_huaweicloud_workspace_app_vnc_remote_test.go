package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppVncRemote_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_app_vnc_remote.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppVncRemote_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "url", regexp.MustCompile(`^https?://.*$`)),
					resource.TestCheckOutput("is_type_valid", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppVncRemote_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_vnc_remote" "test" {
  server_id = "%[1]s"
}

output "is_type_valid" {
  value = strcontains(data.huaweicloud_workspace_app_vnc_remote.test.type, "vnc")
}
`, acceptance.HW_WORKSPACE_APP_SERVER_ID)
}
