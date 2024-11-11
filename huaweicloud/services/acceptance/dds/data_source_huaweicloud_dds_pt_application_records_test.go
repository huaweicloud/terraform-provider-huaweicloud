package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsPtApplicationRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_pt_application_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsPtApplicationRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.applied_at"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.apply_result"),
				),
			},
		},
	})
}

func testDataSourceDdsPtApplicationRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_pt_application_records" "test" {
  depends_on = [huaweicloud_dds_parameter_template_apply.test]

  configuration_id = huaweicloud_dds_parameter_template.test.id
}
`, testParameterTemplateApply_basic(name))
}
