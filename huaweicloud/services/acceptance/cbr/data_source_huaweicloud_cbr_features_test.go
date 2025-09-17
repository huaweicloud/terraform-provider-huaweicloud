package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// The API used by this datasource is currently in public beta and is temporarily unavailable in some regions.
func TestAccDataCbrFeatures_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbr_features.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCbrFeatures_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "feature_value"),
				),
			},
		},
	})
}

const testAccDataCbrFeatures_basic = `
data "huaweicloud_cbr_features" "test" {}
`
