package tms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_key"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
				),
			},
		},
	})
}

func testDataSourceTmsQuotas_basic() string {
	return `
data "huaweicloud_tms_quotas" "test" {}
`
}
