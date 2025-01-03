package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussPtModifyRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_pt_modify_records.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussParameterTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussPtModifyRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.update_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussPtModifyRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_pt_modify_records" "test" {
  config_id = "%s"
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_PARAMETER_TEMPLATE_ID)
}
