package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsPtModificationRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_pt_modification_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_basic(rName),
			},
			{
				Config: testDdsParameterTemplate_basic_update(rName),
			},
			{
				Config: testDataSourceDdsPtModificationRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceDdsPtModificationRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_pt_modification_records" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
}
`, testDdsParameterTemplate_basic_update(name))
}
