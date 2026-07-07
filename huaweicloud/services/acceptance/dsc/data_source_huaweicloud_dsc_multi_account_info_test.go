package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMultiAccountInfo_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_multi_account_info.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMultiAccountInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "is_admin"),
					resource.TestCheckResourceAttrSet(dataSource, "is_delegated_admin"),
					resource.TestCheckResourceAttrSet(dataSource, "is_trusted_service"),
				),
			},
		},
	})
}

const testDataSourceMultiAccountInfo_basic = `
data "huaweicloud_dsc_multi_account_info" "test" {}
`
