package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScriptOrderStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_script_order_statistics.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScriptOrderStatistics_basic(rName, acceptance.HW_COC_INSTANCE_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "execute_statistics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_statistics.0.instance_count"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_statistics.0.instance_status"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScriptOrderStatistics_basic(name, instanceId string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_script_order_statistics" "test" {
  execute_uuid = huaweicloud_coc_script_execute.test.id
}
`, tesScriptExecute_basic(name, instanceId))
}
