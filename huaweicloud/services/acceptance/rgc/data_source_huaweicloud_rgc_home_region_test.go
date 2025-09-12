package rgc

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHomeRegion_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_home_region.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	pattern := `^[a-z]+-[a-z]+-\d$`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHomeRegion_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchOutput("home_region", regexp.MustCompile(pattern)),
					resource.TestCheckResourceAttrSet(dataSource, "home_region"),
				),
			},
		},
	})
}

const testDataSourceHomeRegion_basic = `
data "huaweicloud_rgc_home_region" "test" {}

output "home_region" {
  value = data.huaweicloud_rgc_home_region.test.home_region
}
`
