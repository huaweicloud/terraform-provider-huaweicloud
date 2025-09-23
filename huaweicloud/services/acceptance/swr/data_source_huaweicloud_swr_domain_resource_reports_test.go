package swr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrDomainReports_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_domain_resource_reports.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSwrDomainReports_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "reports.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.date"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.value"),
				),
			},
		},
	})
}

const testDataSourceDataSourceSwrDomainReports_basic = `
data "huaweicloud_swr_domain_resource_reports" "test" {
  resource_type = "store"
  frequency     = "daily"
}
`
