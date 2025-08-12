package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocComponentDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_soc_component_detail.test"
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
				Config: testDataSourceSocComponentDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dev_language"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dev_language_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.alliance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.logo"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.label"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operate_history.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.component_type"),
				),
			},
		},
	})
}

func testDataSourceSocComponentDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_components" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_secmaster_soc_component_detail" "test" {
  workspace_id = "%[1]s"
  component_id = data.huaweicloud_secmaster_soc_components.test.data[0].id
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
