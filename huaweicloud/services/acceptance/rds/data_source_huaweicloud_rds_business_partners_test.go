package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsBusinessPartners_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_business_partners.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsBusinessPartners_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "business_partners.#"),
					resource.TestCheckResourceAttrSet(dataSource, "business_partners.0.order"),
					resource.TestCheckResourceAttrSet(dataSource, "business_partners.0.bp_name"),
					resource.TestCheckResourceAttrSet(dataSource, "business_partners.0.bp_domain_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsBusinessPartners_basic() string {
	return `
data "huaweicloud_rds_business_partners" "test" {}
`
}
