package antiddos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntiDdosQuota_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_antiddos_quota.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAntiDdosQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "current"),
					resource.TestCheckResourceAttrSet(dataSource, "quota"),
				),
			},
		},
	})
}

const testDataSourceAntiDdosQuota_basic = `
data "huaweicloud_antiddos_quota" "test" {}
`
