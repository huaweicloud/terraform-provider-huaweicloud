package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlPtModifyRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_pt_modify_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlPtModifyRecords_base(rName),
			},
			{
				Config: testDataSourceGaussdbMysqlPtModifyRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "histories.#", "2"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.update_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.is_applied"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.applied"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.update_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.is_applied"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.1.applied"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlPtModifyRecords_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name = "%s"
}
`, name)
}

func testDataSourceGaussdbMysqlPtModifyRecords_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name = "%s"

  parameter_values = {
    auto_increment_increment = "4"
    auto_increment_offset    = "5"
  }
}

data "huaweicloud_gaussdb_mysql_pt_modify_records" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_parameter_template.test]

  configuration_id = huaweicloud_gaussdb_mysql_parameter_template.test.id
}
`, name)
}
