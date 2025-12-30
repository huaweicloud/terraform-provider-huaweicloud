package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponentDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_component_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterComponentId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComponentDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "component_id_attr"),
					resource.TestCheckResourceAttrSet(dataSource, "component_name"),
					resource.TestCheckResourceAttrSet(dataSource, "description"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade"),
					resource.TestCheckResourceAttrSet(dataSource, "version"),
				),
			},
		},
	})
}

func testAccDataSourceComponentDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_detail" "test" {
  workspace_id = "%[1]s"
  component_id = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID)
}
