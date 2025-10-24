package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAgencies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_csms_agencies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, ensure that there is a corresponding secret agency.
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgencies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "agency_granted"),
				),
			},
		},
	})
}

const testAccDataSourceAgencies_basic = `
data "huaweicloud_csms_agencies" "test" {
  secret_type = "RDS-FG"
}
`
