package dws

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.limit"),
					resource.TestCheckOutput("is_used_quotas", "true"),
				),
			},
		},
	})
}

const testDataSourceQuotas_basic = `
data "huaweicloud_dws_quotas" "test" {}

output "is_used_quotas" {
  value = length([for v in data.huaweicloud_dws_quotas.test.quotas[*].used : v if v > 0]) > 0
}
`
