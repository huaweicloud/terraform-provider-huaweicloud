package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceParameterModifyRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_instance_parameter_modify_records.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstanceParameterModifyRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.history_id"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.status"),
				),
			},
		},
	})
}

func testDataSourceInstanceParameterModifyRecords_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_parameter_modify_records" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testAccDcsV1Instance_basic(name))
}
