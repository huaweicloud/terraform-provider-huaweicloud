package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBInstanceParametersHistories_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_instance_parameters_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBInstanceParametersHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.update_result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.applied"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.applied_at"),

					resource.TestCheckOutput("param_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBInstanceParametersHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_instance_parameters_histories" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_geminidb_instance_parameters_histories" "param_name_filter" {
  instance_id    = "%[1]s"
  parameter_name = "AuthFailLockTime"
}

output "param_name_filter_is_useful" {
  value = length(data.huaweicloud_geminidb_instance_parameters_histories.param_name_filter.histories) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
