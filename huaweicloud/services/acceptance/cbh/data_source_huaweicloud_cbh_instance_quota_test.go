package cbh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceQuota_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_cbh_instance_quota.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "quota_used"),
				),
			},
		},
	})
}

const testAccDataSourceInstanceQuota_basic string = `
data "huaweicloud_cbh_instance_quota" "test" {}
`
