package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsInstanceParameterModificationRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_instance_parameter_modification_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsInstanceParameterModificationRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.parameter_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.new_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.update_result"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.applied"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.applied_at"),
				),
			},
		},
	})
}

func testDataSourceDdsInstanceParameterModificationRecords_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_instance_parameter_modification_records" "test" {
  depends_on = [huaweicloud_dds_instance_parameters_modify.test]
  
  instance_id = huaweicloud_dds_instance.instance.id
}
`, testAccDDSInstanceV3ModifyParams_basic(name))
}
