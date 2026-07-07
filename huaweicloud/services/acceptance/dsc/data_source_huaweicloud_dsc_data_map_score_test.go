package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDataMapScore_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_data_map_score.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDataMapScore_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "analysis_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_analyze_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "score"),
				),
			},
		},
	})
}

const testAccDataSourceDscDataMapScore_basic = `
data "huaweicloud_dsc_data_map_score" "test" {}
`
