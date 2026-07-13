package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscShowMetadata_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_show_metadata.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscShowMetadata_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "column_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "file_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sensitive_column_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sensitive_file_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_num"),
				),
			},
		},
	})
}

const testAccDataSourceDscShowMetadata_basic = `
data "huaweicloud_dsc_show_metadata" "test" {}
`
