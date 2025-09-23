package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponents_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_components.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceComponents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.component_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.upgrade"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.maintainer"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.description"),
				),
			},
		},
	})
}

func testDataSourceComponents_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_components" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
