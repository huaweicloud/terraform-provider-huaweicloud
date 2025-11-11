package fgs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please ensure that you have created trusted agencies for FGS service.
func TestAccDataServiceTrustedAgencies_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_service_trusted_agencies.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataServiceTrustedAgencies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "agencies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "agencies.0.name"),
					// When the agency never expires, the "expire_time" is empty string.
				),
			},
		},
	})
}

const testAccDataServiceTrustedAgencies_basic = `
data "huaweicloud_fgs_service_trusted_agencies" "test" {}
`
