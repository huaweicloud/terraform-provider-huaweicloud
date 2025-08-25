package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case cannot be executed successfully.
func TestAccDataSourceFlowBlock_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_aad_flow_block.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFlowBlock_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ips.#"),
				),
			},
		},
	})
}

func testDataSourceFlowBlock_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_flow_block" "test" {
  instance_id = "%s"
}
`, acceptance.HW_AAD_INSTANCE_ID)
}
