package cdn

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBillingOption_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_billing_option.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		allFilterRName = "data.huaweicloud_cdn_billing_option.all_filter"
		allFilterDc    = acceptance.InitDataSourceCheck(allFilterRName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBillingOption_expectError,
				ExpectError: regexp.MustCompile("Your query returned no results. " +
					"Please change your search criteria and try again."),
			},
			{
				Config: testAccDataBillingOption_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "product_type"),
					resource.TestCheckResourceAttrSet(rName, "service_area"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "charge_mode"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "effective_time"),

					allFilterDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(allFilterRName, "product_type"),
					resource.TestCheckResourceAttrSet(allFilterRName, "service_area"),
					resource.TestCheckResourceAttrSet(allFilterRName, "status"),
					resource.TestCheckResourceAttrSet(allFilterRName, "charge_mode"),
					resource.TestCheckResourceAttrSet(allFilterRName, "created_at"),
					resource.TestCheckResourceAttrSet(allFilterRName, "effective_time"),
				),
			},
		},
	})
}

const testAccDataBillingOption_expectError = `
data "huaweicloud_cdn_billing_option" "test" {
  product_type = "base"
}

data "huaweicloud_cdn_billing_option" "expect_error" {
  product_type = "base"
  status       = data.huaweicloud_cdn_billing_option.test.status == "active" ? "upcoming" : "active"
}
`

const testAccDataBillingOption_basic = `
data "huaweicloud_cdn_billing_option" "test" {
  product_type = "base"
}

data "huaweicloud_cdn_billing_option" "all_filter" {
  product_type = "base"
  status       = data.huaweicloud_cdn_billing_option.test.status
  service_area = data.huaweicloud_cdn_billing_option.test.service_area
}
`
