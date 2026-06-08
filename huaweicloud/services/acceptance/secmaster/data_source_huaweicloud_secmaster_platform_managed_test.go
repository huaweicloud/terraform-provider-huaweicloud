package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlatformManaged_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_platform_managed.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePlatformManaged_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "dw_region"),
					resource.TestCheckResourceAttrSet(dataSource, "publish_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tenant_managed_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "whitelist"),
				),
			},
		},
	})
}

func testDataSourcePlatformManaged_basic() string {
	return `
data "huaweicloud_secmaster_platform_managed" "test" {}
`
}
