package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocPublicScriptDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_public_script_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocScriptID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocPublicScriptDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "description"),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "content"),
					resource.TestCheckResourceAttrSet(dataSource, "script_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "script_params.0.param_name"),
					resource.TestCheckResourceAttrSet(dataSource, "script_params.0.param_description"),
					resource.TestCheckResourceAttrSet(dataSource, "script_params.0.param_order"),
					resource.TestCheckResourceAttrSet(dataSource, "script_params.0.sensitive"),
					resource.TestCheckResourceAttrSet(dataSource, "gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "properties.#"),
					resource.TestCheckResourceAttrSet(dataSource, "properties.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "properties.0.version"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocPublicScriptDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_public_script_detail" "test" {
  script_uuid = "%s"
}
`, acceptance.HW_COC_SCRIPT_ID)
}
