package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityDatasourcePermissions_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_security_datasource_permissions.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityDatasourcePermissions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "datasource_permissions.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "datasource_permissions.0.datasource_type"),
					resource.TestMatchResourceAttr(all, "datasource_permissions.0.permission_types.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestMatchResourceAttr(all, "datasource_permissions.0.permission_actions.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
				),
			},
		},
	})
}

func testAccDataSecurityDatasourcePermissions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_security_datasource_permissions" "test" {
  workspace_id = "%s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
