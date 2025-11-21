package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceControlDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_control_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCControlID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceControlDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "aliases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.content.#"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.content.0.ch"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.content.0.en"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "behavior"),
					resource.TestCheckResourceAttrSet(dataSource, "control_objective"),
					resource.TestCheckResourceAttrSet(dataSource, "framework.#"),
					resource.TestCheckResourceAttrSet(dataSource, "guidance"),
					resource.TestCheckResourceAttrSet(dataSource, "identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "implementation"),
					resource.TestCheckResourceAttrSet(dataSource, "owner"),
					resource.TestCheckResourceAttrSet(dataSource, "release_date"),
					resource.TestCheckResourceAttrSet(dataSource, "resource.#"),
					resource.TestCheckResourceAttrSet(dataSource, "service"),
					resource.TestCheckResourceAttrSet(dataSource, "severity"),
					resource.TestCheckResourceAttrSet(dataSource, "version"),
				),
			},
		},
	})
}

func testDataSourceDataSourceControlDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_control_detail" "test" {
  control_id = "%[1]s"
}
`, acceptance.HW_RGC_CONTROL_ID)
}
