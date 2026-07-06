package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDataMapSecurityLevel_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_data_map_security_level.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDataMapSecurityLevel_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "level"),
				),
			},
		},
	})
}

const testAccDataSourceDscDataMapSecurityLevel_basic = `
data "huaweicloud_dsc_data_map_security_level" "test" {}
`
