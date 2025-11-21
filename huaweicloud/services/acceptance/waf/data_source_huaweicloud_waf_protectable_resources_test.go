package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafProtectables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_waf_protectable_resources.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafProtectables_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
				),
			},
		},
	})
}

const testDataSourceDataSourceWafProtectables_basic = `
data "huaweicloud_waf_protectable_resources" "test" {
  resource_type = "elb"
}
`
