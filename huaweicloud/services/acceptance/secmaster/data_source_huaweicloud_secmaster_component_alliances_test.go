package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponentAlliances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_component_alliances.test"
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
				Config: testDataSourceComponentAlliances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.alliance_code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.alliance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.alliance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.alliance_description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.logo"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
				),
			},
		},
	})
}

func testDataSourceComponentAlliances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_component_alliances" "test" {
  workspace_id = "%[1]s"
  is_built_in  = false
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
