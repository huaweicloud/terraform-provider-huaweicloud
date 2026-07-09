package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscProtectDataTypes_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_protect_data_types.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscProtectDataTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.protect_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.life_cycle"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.update_time"),
				),
			},
		},
	})
}

const testAccDataSourceDscProtectDataTypes_basic = `
data "huaweicloud_dsc_protect_data_types" "test" {
  life_cycle = "TRANSMISSION"
}
`
