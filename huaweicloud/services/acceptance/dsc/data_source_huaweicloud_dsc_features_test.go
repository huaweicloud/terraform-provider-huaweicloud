package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscFeatures_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_features.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscFeatures_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.name"),
				),
			},
		},
	})
}

const testAccDataSourceDscFeatures_basic = `
data "huaweicloud_dsc_features" "test" {}
`
