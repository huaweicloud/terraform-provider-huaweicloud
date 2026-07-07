package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDbssOemInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_dbss_oem_info.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDbssOemInfo_conf,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "ins_info_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscDbssOemInfo_conf = `
data "huaweicloud_dsc_dbss_oem_info" "test" {
  type = "DBSS"
}
`
