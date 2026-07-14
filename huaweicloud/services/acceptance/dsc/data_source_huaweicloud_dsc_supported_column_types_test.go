package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSupportedColumnTypes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_supported_column_types.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSupportedColumnTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_source_type_response"),
					resource.TestCheckResourceAttrSet(dataSource, "supported_column_types.#"),
				),
			},
		},
	})
}

const testDataSourceSupportedColumnTypes_basic = `
data "huaweicloud_dsc_supported_column_types" "test" {
  data_source_type = "MRS_HIVE"
}
`
